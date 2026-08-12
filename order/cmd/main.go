package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiv1 "github.com.Saveliy12/microservices2.git/order/internal/api/order/v1"
	"github.com.Saveliy12/microservices2.git/order/internal/client/db"
	dbpg "github.com.Saveliy12/microservices2.git/order/internal/client/db/pg"
	inventoryclient "github.com.Saveliy12/microservices2.git/order/internal/clients/inventory"
	paymentclient "github.com.Saveliy12/microservices2.git/order/internal/clients/payment"
	"github.com.Saveliy12/microservices2.git/order/internal/migrator"
	repositoryorder "github.com.Saveliy12/microservices2.git/order/internal/repository/order"
	service "github.com.Saveliy12/microservices2.git/order/internal/service"
	serviceorder "github.com.Saveliy12/microservices2.git/order/internal/service/order"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort             = ":8080"
	readHeaderTimeout    = 5 * time.Second
	shutdownTimeout      = 10 * time.Second
	swaggerJSONPath      = "shared/api/bundles/order.openapi.v1.bundle.yaml"
	swaggerStaticDir     = "shared/api"
	defaultMigrationsDir = "order/migrations"
	defaultInventoryAddr = "localhost:50051"
	defaultPaymentAddr   = "localhost:50052"
)

// inventoryAddr возвращает адрес Inventory-сервиса из окружения
// с дефолтом для локального запуска
func inventoryAddr() string {
	if addr := os.Getenv("INVENTORY_ADDR"); addr != "" {
		return addr
	}

	return defaultInventoryAddr
}

// paymentAddr возвращает адрес Payment-сервиса из окружения
// с дефолтом для локального запуска
func paymentAddr() string {
	if addr := os.Getenv("PAYMENT_ADDR"); addr != "" {
		return addr
	}

	return defaultPaymentAddr
}

// initializeConnections устанавливает gRPC соединения с микросервисами инвентаря и платежей
// Если соединение с платежами упадёт, соединение с инвентарём будет закрыто и вернётся ошибка
func initializeConnections(ctx context.Context) (*grpc.ClientConn, *grpc.ClientConn, error) {
	invConn, err := connectGRPC(ctx, inventoryAddr())
	if err != nil {
		return nil, nil, err
	}

	payConn, err := connectGRPC(ctx, paymentAddr())
	if err != nil {
		_ = invConn.Close()
		return nil, nil, err
	}

	return invConn, payConn, nil
}

// buildServiceStack собирает слои чистой архитектуры: репозиторий -> сервис -> API
func buildServiceStack(dbClient db.Client, invClient service.InventoryClient, payClient service.PaymentClient) *apiv1.API {
	repo := repositoryorder.NewRepository(dbClient)
	service := serviceorder.NewService(repo, invClient, payClient)
	return apiv1.NewAPI(service)
}

// registerRoutes регистрирует все HTTP обработчики для API эндпоинтов заказов
func registerRoutes(mux *http.ServeMux, api *apiv1.API) {
	fileServer := http.FileServer(http.Dir(swaggerStaticDir))

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

	mux.HandleFunc("/swagger.json", serveSwaggerJSON)
	mux.Handle("/swagger-ui.html", fileServer)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
	})
}

// deferCloseConn безопасно закрывает gRPC соединение и логирует любые ошибки
func deferCloseConn(conn *grpc.ClientConn, name string) {
	if err := conn.Close(); err != nil {
		log.Printf("failed to close %s connection: %v", name, err)
	}
}

// initDB создаёт пул соединений PostgreSQL и применяет миграции
// Возвращает обёртку над пулом (db.Client) и функцию закрытия
func initDB(ctx context.Context) (db.Client, func(), error) {
	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		return nil, nil, errors.New("DB_URI is not set")
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = defaultMigrationsDir
	}
	mig := migrator.NewMigrator(sqlDB, migrationsDir)
	if err := mig.Up(); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return dbpg.NewPGClient(pool), pool.Close, nil
}

func main() {
	// Загружаем переменные окружения из .env (если он есть) перед чтением настроек
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	// Настройка graceful shutdown с обработкой сигналов
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инициализация соединений с микросервисами
	invConn, payConn, err := initializeConnections(ctx)
	if err != nil {
		stop()
		log.Fatalf("failed to initialize connections: %v", err)
	}
	defer deferCloseConn(invConn, "inventory")
	defer deferCloseConn(payConn, "payment")

	// Инициализация пула соединений с базой данных и применение миграций
	dbClient, closeDB, err := initDB(ctx)
	if err != nil {
		stop()
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer closeDB()

	api := buildServiceStack(dbClient, inventoryclient.New(invConn), paymentclient.New(payConn))

	// Регистрация HTTP маршрутов и обработчиков
	mux := http.NewServeMux()
	registerRoutes(mux, api)

	// Конфигурация и запуск HTTP сервера
	server := &http.Server{
		Addr:              httpPort,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 Order HTTP server listening on %s", httpPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start Order HTTP server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	log.Println("🛑 Shutting down Order HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("failed to shutdown Order HTTP server: %v", err)
	}

	log.Println("✅ Order server stopped")
}

// connectGRPC устанавливает gRPC соединение с блокировкой до готовности
func connectGRPC(ctx context.Context, address string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// serveSwaggerJSON отдаёт OpenAPI спецификацию в формате JSON
func serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, swaggerJSONPath)
}
