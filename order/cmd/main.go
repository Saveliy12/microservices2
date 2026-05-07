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
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	swaggerJSONPath   = "shared/api/bundles/order.openapi.v1.bundle.yaml"
	swaggerStaticDir  = "shared/api"
)

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

// Adapter for payment gRPC client
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

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inventoryConn, err := connectGRPC(ctx, "localhost:50051")
	if err != nil {
		stop()
		log.Fatalf("failed to connect to inventory: %v", err)
	}

	paymentConn, err := connectGRPC(ctx, "localhost:50052")
	if err != nil {
		// ensure we stop the signal context to release resources
		_ = inventoryConn.Close()
		stop()
		log.Fatalf("failed to connect to payment: %v", err)
	}

	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("failed to close inventory connection: %v", err)
		}
	}()

	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v", err)
		}
	}()

	orderRepository := repositoryorder.NewRepository()

	inventoryClient := &inventoryClientAdapter{
		client: inventoryv1.NewInventoryServiceClient(inventoryConn),
	}

	paymentClient := &paymentClientAdapter{
		client: paymentv1.NewPaymentServiceClient(paymentConn),
	}

	orderService := serviceorder.NewService(orderRepository, inventoryClient, paymentClient)

	api := apiv1.NewAPI(orderService)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(swaggerStaticDir))

	// Register create endpoint for exact path
	mux.HandleFunc("/api/v1/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.CreateOrder(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	// Register prefix for operations on specific order resources
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

	<-ctx.Done()
	log.Println("🛑 Shutting down Order HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("failed to shutdown Order HTTP server: %v", err)
	}

	log.Println("✅ Order server stopped")
}

func connectGRPC(ctx context.Context, address string) (*grpc.ClientConn, error) {
	// DialContext with WithBlock will wait until the connection is ready or context is done
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, swaggerJSONPath)
}
