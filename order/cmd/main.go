package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiv1 "github.com.Saveliy12/microservices2.git/order/internal/api/order/v1"
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repositoryorder "github.com.Saveliy12/microservices2.git/order/internal/repository/order"
	serviceorder "github.com.Saveliy12/microservices2.git/order/internal/service/order"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	httpPort          = ":8080"
	inventoryAddr     = "localhost:50051"
	paymentAddr       = "localhost:50052"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	swaggerJSONPath   = "shared/api/bundles/order.openapi.v1.bundle.yaml"
	swaggerStaticDir  = "shared/api"
)

// inventoryClientAdapter реализует интерфейс service.InventoryClient для gRPC взаимодействия с микросервисом инвентаря
type inventoryClientAdapter struct {
	client inventoryv1.InventoryServiceClient
}

func (a *inventoryClientAdapter) ListParts(ctx context.Context, partUUIDs []string) ([]model.PartQuote, error) {
	req := &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: make([]*wrapperspb.StringValue, 0, len(partUUIDs)),
		},
	}
	for _, uuid := range partUUIDs {
		req.Filter.Uuids = append(req.Filter.Uuids, wrapperspb.String(uuid))
	}

	resp, err := a.client.ListParts(ctx, req)
	if err != nil {
		return nil, err
	}

	parts := make([]model.PartQuote, 0, len(resp.GetParts()))
	for _, part := range resp.GetParts() {
		parts = append(parts, model.PartQuote{
			UUID:  part.GetUuid(),
			Price: part.GetPrice(),
		})
	}

	return parts, nil
}

// paymentClientAdapter реализует интерфейс service.PaymentClient для gRPC взаимодействия с микросервисом платежей
type paymentClientAdapter struct {
	client paymentv1.PaymentServiceClient
}

func (a *paymentClientAdapter) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	pmEnum, err := paymentMethodToProto(paymentMethod)
	if err != nil {
		return "", err
	}

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: pmEnum,
	}

	resp, err := a.client.PayOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetTransactionUuid(), nil
}

// paymentMethodToProto преобразует доменный enum PaymentMethod в protobuf enum PaymentMethod
func paymentMethodToProto(method model.PaymentMethod) (paymentv1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case model.PaymentMethodSbp:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, errors.New("invalid payment method")
	}
}

// initializeConnections устанавливает gRPC соединения с микросервисами инвентаря и платежей
// Если соединение с платежами упадёт, соединение с инвентарём будет закрыто и вернётся ошибка
func initializeConnections(ctx context.Context) (*grpc.ClientConn, *grpc.ClientConn, error) {
	invConn, err := connectGRPC(ctx, inventoryAddr)
	if err != nil {
		return nil, nil, err
	}

	payConn, err := connectGRPC(ctx, paymentAddr)
	if err != nil {
		_ = invConn.Close()
		return nil, nil, err
	}

	return invConn, payConn, nil
}

// newInventoryAdapter создаёт адаптер для взаимодействия с микросервисом инвентаря
func newInventoryAdapter(conn *grpc.ClientConn) *inventoryClientAdapter {
	return &inventoryClientAdapter{
		client: inventoryv1.NewInventoryServiceClient(conn),
	}
}

// newPaymentAdapter создаёт адаптер для взаимодействия с микросервисом платежей
func newPaymentAdapter(conn *grpc.ClientConn) *paymentClientAdapter {
	return &paymentClientAdapter{
		client: paymentv1.NewPaymentServiceClient(conn),
	}
}

// buildServiceStack собирает слои чистой архитектуры: репозиторий -> сервис -> API
func buildServiceStack(invAdapter *inventoryClientAdapter, payAdapter *paymentClientAdapter) *apiv1.API {
	repo := repositoryorder.NewRepository()
	service := serviceorder.NewService(repo, invAdapter, payAdapter)
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

// deferCloseConn безопасно закрывает gRPC соединение и логирует любые ошибки.
func deferCloseConn(conn *grpc.ClientConn, name string) {
	if err := conn.Close(); err != nil {
		log.Printf("failed to close %s connection: %v", name, err)
	}
}

func main() {
	// Настройка graceful shutdown с обработкой сигналов.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инициализация соединений с микросервисами.
	invConn, payConn, err := initializeConnections(ctx)
	if err != nil {
		stop()
		log.Fatalf("failed to initialize connections: %v", err)
	}
	defer deferCloseConn(invConn, "inventory")
	defer deferCloseConn(payConn, "payment")

	// Построение стека чистой архитектуры с адаптерами.
	api := buildServiceStack(newInventoryAdapter(invConn), newPaymentAdapter(payConn))

	// Регистрация HTTP маршрутов и обработчиков.
	mux := http.NewServeMux()
	registerRoutes(mux, api)

	// Конфигурация и запуск HTTP сервера.
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

	// Ожидание сигнала завершения.
	<-ctx.Done()
	log.Println("🛑 Shutting down Order HTTP server...")

	// Graceful shutdown с таймаутом.
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

// serveSwaggerJSON отдаёт OpenAPI спецификацию в формате JSON.
func serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, swaggerJSONPath)
}
