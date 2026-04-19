package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Status string

type PaymentMethod string

const (
	httpPort          = ":8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	swaggerJSONPath   = "shared/api/bundles/order.openapi.v1.bundle.yaml"
	swaggerStaticDir  = "shared/api"

	StatusPendingPayment Status = "PENDING_PAYMENT"
	StatusPaid           Status = "PAID"
	StatusCancelled      Status = "CANCELLED"
)

type Order struct {
	OrderUUID       string         `json:"order_uuid"`
	UserUUID        string         `json:"user_uuid"`
	PartUUIDs       []string       `json:"part_uuids"`
	TotalPrice      float64        `json:"total_price"`
	TransactionUUID *string        `json:"transaction_uuid,omitempty"`
	PaymentMethod   *PaymentMethod `json:"payment_method,omitempty"`
	Status          Status         `json:"status"`
}

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{orders: make(map[string]*Order)}
}

func (s *OrderStorage) Create(order *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.OrderUUID] = order
}

func (s *OrderStorage) Get(orderUUID string) (*Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[orderUUID]
	return order, ok
}

func (s *OrderStorage) Update(order *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.OrderUUID] = order
}

type OrderHandler struct {
	storage       *OrderStorage
	inventory     inventoryv1.InventoryServiceClient
	payment       paymentv1.PaymentServiceClient
	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn
}

func NewOrderHandler(storage *OrderStorage, inventory inventoryv1.InventoryServiceClient, payment paymentv1.PaymentServiceClient, inventoryConn *grpc.ClientConn, paymentConn *grpc.ClientConn) *OrderHandler {
	return &OrderHandler{
		storage:       storage,
		inventory:     inventory,
		payment:       payment,
		inventoryConn: inventoryConn,
		paymentConn:   paymentConn,
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inventoryConn, err := connectGRPC(ctx, "localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to inventory: %v", err)
	}
	defer inventoryConn.Close()

	paymentConn, err := connectGRPC(ctx, "localhost:50052")
	if err != nil {
		log.Fatalf("failed to connect to payment: %v", err)
	}
	defer paymentConn.Close()

	handler := NewOrderHandler(
		NewOrderStorage(),
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
		inventoryConn,
		paymentConn,
	)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(swaggerStaticDir))
	mux.HandleFunc("POST /api/v1/orders", handler.createOrder)
	mux.HandleFunc("GET /api/v1/orders/", handler.getOrder)
	mux.HandleFunc("POST /api/v1/orders/", handler.orderAction)
	mux.HandleFunc("GET /swagger.json", serveSwaggerJSON)
	mux.Handle("/swagger-ui.html", fileServer)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
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
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	conn.Connect()
	for state := conn.GetState(); state != connectivity.Ready; state = conn.GetState() {
		if !conn.WaitForStateChange(ctx, state) {
			_ = conn.Close()
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, context.DeadlineExceeded
		}
	}

	return conn, nil
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {

	type request struct {
		UserUUID  string   `json:"user_uuid"`
		PartUUIDs []string `json:"part_uuids"`
	}

	type response struct {
		OrderUUID  string  `json:"order_uuid"`
		TotalPrice float64 `json:"total_price"`
	}

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.UserUUID == "" || len(req.PartUUIDs) == 0 {
		writeError(w, http.StatusBadRequest, "user_uuid and part_uuids are required")
		return
	}

	partsRequest := &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: make([]*wrapperspb.StringValue, 0, len(req.PartUUIDs)),
		},
	}

	for _, partUUID := range req.PartUUIDs {
		partsRequest.Filter.Uuids = append(partsRequest.Filter.Uuids, wrapperspb.String(partUUID))
	}

	partsResponse, err := h.inventory.ListParts(r.Context(), partsRequest)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to query inventory")
		return
	}

	partsByUUID := make(map[string]*inventoryv1.Part, len(partsResponse.GetParts()))
	for _, part := range partsResponse.GetParts() {
		partsByUUID[part.GetUuid()] = part
	}

	totalPrice := 0.0
	for _, partUUID := range req.PartUUIDs {
		part, ok := partsByUUID[partUUID]
		if !ok {
			writeError(w, http.StatusBadRequest, "one or more parts not found")
			return
		}
		totalPrice += part.GetPrice()
	}

	order := &Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   req.UserUUID,
		PartUUIDs:  req.PartUUIDs,
		TotalPrice: totalPrice,
		Status:     StatusPendingPayment,
	}

	h.storage.Create(order)
	writeJSON(w, http.StatusCreated, response{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	})
}

func (h *OrderHandler) getOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	orderUUID, err := extractOrderUUID(r.URL.Path, "/api/v1/orders/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order uuid")
		return
	}

	order, ok := h.storage.Get(orderUUID)
	if !ok {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) orderAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	switch {
	case strings.HasSuffix(r.URL.Path, "/pay"):
		h.payOrder(w, r)
	case strings.HasSuffix(r.URL.Path, "/cancel"):
		h.cancelOrder(w, r)
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}

func (h *OrderHandler) payOrder(w http.ResponseWriter, r *http.Request) {
	type request struct {
		PaymentMethod PaymentMethod `json:"payment_method"`
	}
	type response struct {
		TransactionUUID string `json:"transaction_uuid"`
	}

	orderUUID, err := extractActionOrderUUID(r.URL.Path, "/pay")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order uuid")
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.PaymentMethod == "" {
		writeError(w, http.StatusBadRequest, "payment_method is required")
		return
	}

	order, ok := h.storage.Get(orderUUID)
	if !ok {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}
	if order.Status == StatusPaid {
		writeError(w, http.StatusConflict, "order already paid")
		return
	}
	if order.Status == StatusCancelled {
		writeError(w, http.StatusConflict, "order is cancelled")
		return
	}

	paymentMethod, paymentMethodEnum, err := parsePaymentMethod(string(req.PaymentMethod))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	paymentResponse, err := h.payment.PayOrder(r.Context(), &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentMethodEnum,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to call payment service")
		return
	}

	transactionUUID := paymentResponse.GetTransactionUuid()
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &paymentMethod
	order.Status = StatusPaid
	h.storage.Update(order)

	writeJSON(w, http.StatusOK, response{TransactionUUID: transactionUUID})
}

func (h *OrderHandler) cancelOrder(w http.ResponseWriter, r *http.Request) {
	orderUUID, err := extractActionOrderUUID(r.URL.Path, "/cancel")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order uuid")
		return
	}

	order, ok := h.storage.Get(orderUUID)
	if !ok {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}
	if order.Status == StatusPaid {
		writeError(w, http.StatusConflict, "paid order cannot be cancelled")
		return
	}

	order.Status = StatusCancelled
	h.storage.Update(order)
	w.WriteHeader(http.StatusNoContent)
}

func extractOrderUUID(path, prefix string) (string, error) {
	if !strings.HasPrefix(path, prefix) {
		return "", errors.New("invalid path")
	}

	value := strings.TrimPrefix(path, prefix)
	value = strings.Trim(value, "/")
	if value == "" || strings.Contains(value, "/") {
		return "", errors.New("invalid order uuid")
	}

	return value, nil
}

func extractActionOrderUUID(path, suffix string) (string, error) {
	if !strings.HasSuffix(path, suffix) {
		return "", errors.New("invalid path")
	}

	value := strings.TrimSuffix(path, suffix)
	value = strings.TrimSuffix(value, "/")
	return extractOrderUUID(value, "/api/v1/orders/")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parsePaymentMethod(value string) (PaymentMethod, paymentv1.PaymentMethod, error) {
	switch value {
	case "CARD", "PAYMENT_METHOD_CARD":
		return PaymentMethod("CARD"), paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case "SBP", "PAYMENT_METHOD_SBP":
		return PaymentMethod("SBP"), paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case "CREDIT_CARD", "PAYMENT_METHOD_CREDIT_CARD":
		return PaymentMethod("CREDIT_CARD"), paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case "INVESTOR_MONEY", "PAYMENT_METHOD_INVESTOR_MONEY":
		return PaymentMethod("INVESTOR_MONEY"), paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return "", paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, errors.New("invalid payment_method")
	}
}

func serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	http.ServeFile(w, r, swaggerJSONPath)
}
