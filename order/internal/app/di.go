package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com.Saveliy12/microservices2.git/order/internal/api/health"
	apiv1 "github.com.Saveliy12/microservices2.git/order/internal/api/order/v1"
	"github.com.Saveliy12/microservices2.git/order/internal/client/db"
	dbpg "github.com.Saveliy12/microservices2.git/order/internal/client/db/pg"
	inventoryclient "github.com.Saveliy12/microservices2.git/order/internal/clients/inventory"
	paymentclient "github.com.Saveliy12/microservices2.git/order/internal/clients/payment"
	"github.com.Saveliy12/microservices2.git/order/internal/config"
	"github.com.Saveliy12/microservices2.git/order/internal/repository"
	repositoryorder "github.com.Saveliy12/microservices2.git/order/internal/repository/order"
	"github.com.Saveliy12/microservices2.git/order/internal/service"
	serviceorder "github.com.Saveliy12/microservices2.git/order/internal/service/order"
	"github.com.Saveliy12/microservices2.git/platform/pkg/closer"
	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
	migratorpg "github.com.Saveliy12/microservices2.git/platform/pkg/migrator/pg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	config          config.OrderConfig
	pgPool          *pgxpool.Pool
	dbClient        db.Client
	inventoryConn   *grpc.ClientConn
	paymentConn     *grpc.ClientConn
	inventoryClient service.InventoryClient
	paymentClient   service.PaymentClient
	orderRepository repository.OrderRepository
	orderService    service.OrderService
	orderAPI        *apiv1.API
	httpServer      *http.Server
	closer          *closer.Closer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) Config() (config.OrderConfig, error) {
	if sp.config == nil {
		cfg, err := config.New()
		if err != nil {
			return nil, err
		}

		sp.config = cfg
	}

	return sp.config, nil
}

// PGPool создаёт пул соединений Postgres
func (sp *serviceProvider) PGPool(ctx context.Context) (*pgxpool.Pool, error) {
	if sp.pgPool == nil {
		cfg, err := sp.Config()
		if err != nil {
			return nil, err
		}

		pool, err := pgxpool.New(ctx, cfg.Postgres().URI())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres: %w", err)
		}

		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			return nil, fmt.Errorf("failed to ping postgres: %w", err)
		}

		sp.pgPool = pool
		sp.Closer().AddNamed("pg_pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})
	}

	return sp.pgPool, nil
}

// DBClient обёртка над пулом + миграции
func (sp *serviceProvider) DBClient(ctx context.Context) (db.Client, error) {
	if sp.dbClient == nil {
		pool, err := sp.PGPool(ctx)
		if err != nil {
			return nil, err
		}

		cfg, err := sp.Config()
		if err != nil {
			return nil, err
		}

		sqlDB := stdlib.OpenDBFromPool(pool)
		sp.Closer().AddNamed("sql_db", func(_ context.Context) error {
			return sqlDB.Close()
		})

		mig := migratorpg.New(sqlDB, cfg.Postgres().MigrationsDir())
		if err := mig.Up(); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}

		sp.dbClient = dbpg.NewPGClient(pool)
	}

	return sp.dbClient, nil
}

// connectGRPC устанавливает соединение с соседним сервисом.
// grpc.NewClient создаёт ленивое соединение, поэтому мы явно дожидаемся состояния Ready:
// если соседний сервис недоступен, приложение падает на старте, а не на первом запросе пользователя.
func connectGRPC(ctx context.Context, address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create grpc client for %s: %w", address, err)
	}

	conn.Connect()
	for state := conn.GetState(); state != connectivity.Ready; state = conn.GetState() {
		if !conn.WaitForStateChange(ctx, state) {
			_ = conn.Close()

			return nil, fmt.Errorf("connect to %s: %w", address, ctx.Err())
		}
	}

	return conn, nil
}

func (sp *serviceProvider) InventoryConn(ctx context.Context) (*grpc.ClientConn, error) {
	if sp.inventoryConn == nil {
		cfg, err := sp.Config()
		if err != nil {
			return nil, err
		}

		addr := net.JoinHostPort(cfg.InventoryGRPC().Host(), fmt.Sprintf("%d", cfg.InventoryGRPC().Port()))

		conn, err := connectGRPC(ctx, addr)
		if err != nil {
			return nil, err
		}

		sp.inventoryConn = conn
		sp.Closer().AddNamed("inventory_grpc", func(_ context.Context) error {
			return conn.Close()
		})
	}

	return sp.inventoryConn, nil
}

func (sp *serviceProvider) PaymentConn(ctx context.Context) (*grpc.ClientConn, error) {
	if sp.paymentConn == nil {
		cfg, err := sp.Config()
		if err != nil {
			return nil, err
		}

		addr := net.JoinHostPort(cfg.PaymentGRPC().Host(), fmt.Sprintf("%d", cfg.PaymentGRPC().Port()))

		conn, err := connectGRPC(ctx, addr)
		if err != nil {
			return nil, err
		}

		sp.paymentConn = conn
		sp.Closer().AddNamed("payment_grpc", func(_ context.Context) error {
			return conn.Close()
		})
	}

	return sp.paymentConn, nil
}

func (sp *serviceProvider) InventoryClient(ctx context.Context) (service.InventoryClient, error) {
	if sp.inventoryClient == nil {
		conn, err := sp.InventoryConn(ctx)
		if err != nil {
			return nil, err
		}

		sp.inventoryClient = inventoryclient.New(conn)
	}

	return sp.inventoryClient, nil
}

func (sp *serviceProvider) PaymentClient(ctx context.Context) (service.PaymentClient, error) {
	if sp.paymentClient == nil {
		conn, err := sp.PaymentConn(ctx)
		if err != nil {
			return nil, err
		}

		sp.paymentClient = paymentclient.New(conn)
	}

	return sp.paymentClient, nil
}

func (sp *serviceProvider) OrderRepository(ctx context.Context) (repository.OrderRepository, error) {
	if sp.orderRepository == nil {
		client, err := sp.DBClient(ctx)
		if err != nil {
			return nil, err
		}

		sp.orderRepository = repositoryorder.NewRepository(client)
	}

	return sp.orderRepository, nil
}

func (sp *serviceProvider) OrderService(ctx context.Context) (service.OrderService, error) {
	if sp.orderService == nil {
		repo, err := sp.OrderRepository(ctx)
		if err != nil {
			return nil, err
		}

		invClient, err := sp.InventoryClient(ctx)
		if err != nil {
			return nil, err
		}

		payClient, err := sp.PaymentClient(ctx)
		if err != nil {
			return nil, err
		}

		sp.orderService = serviceorder.NewService(repo, invClient, payClient)
	}

	return sp.orderService, nil
}

func (sp *serviceProvider) OrderAPI(ctx context.Context) (*apiv1.API, error) {
	if sp.orderAPI == nil {
		svc, err := sp.OrderService(ctx)
		if err != nil {
			return nil, err
		}

		sp.orderAPI = apiv1.NewAPI(svc)
	}

	return sp.orderAPI, nil
}

func (sp *serviceProvider) registerRoutes(mux *http.ServeMux, api *apiv1.API) {
	fileServer := http.FileServer(http.Dir("shared/api"))

	mux.Handle("/health", health.NewHandler())

	mux.HandleFunc("/api/v1/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		api.CreateOrder(w, r)
	})

	mux.HandleFunc("/api/v1/orders/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetOrder(w, r)
		case http.MethodPost:
			api.OrderAction(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/swagger.json", sp.serveSwaggerJSON)
	mux.Handle("/swagger-ui.html", fileServer)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
	})
}

func (sp *serviceProvider) serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "shared/api/bundles/order.openapi.v1.bundle.yaml")
}

// HTTPServer собирает http.Server
func (sp *serviceProvider) HTTPServer(ctx context.Context) (*http.Server, error) {
	if sp.httpServer == nil {
		cfg, err := sp.Config()
		if err != nil {
			return nil, err
		}

		api, err := sp.OrderAPI(ctx)
		if err != nil {
			return nil, err
		}

		mux := http.NewServeMux()
		sp.registerRoutes(mux, api)

		addr := net.JoinHostPort(cfg.HTTP().Host(), fmt.Sprintf("%d", cfg.HTTP().Port()))

		sp.httpServer = &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}
	}

	return sp.httpServer, nil
}

func (sp *serviceProvider) Closer() *closer.Closer {
	if sp.closer == nil {
		sp.closer = closer.NewWithLogger(logger.Logger())
	}

	return sp.closer
}
