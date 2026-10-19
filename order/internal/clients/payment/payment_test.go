package payment_test

import (
	"context"
	"net"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/clients/payment"
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	paymentv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type PaymentClientSuite struct {
	suite.Suite
	ctx    context.Context
	conn   *grpc.ClientConn
	client *payment.Client
	server *paymentServiceServer
}

func (s *PaymentClientSuite) SetupTest() {
	s.ctx = context.Background()

	s.server = newPaymentServiceServer()
	lis := bufconn.Listen(1 << 20)

	grpcServer := grpc.NewServer()
	paymentv1.RegisterPaymentServiceServer(grpcServer, s.server)
	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)

	s.conn = conn
	s.client = payment.New(conn)
}

func (s *PaymentClientSuite) TearDownTest() {
	_ = s.conn.Close()
}

func (s *PaymentClientSuite) TestPayOrder_Success() {
	s.server.payOrderResponse = &paymentv1.PayOrderResponse{TransactionUuid: "txn-1"}

	transactionUUID, err := s.client.PayOrder(s.ctx, "order-1", "user-1", model.PaymentMethodCard)

	s.Require().NoError(err)
	s.Equal("txn-1", transactionUUID)
	s.Equal("order-1", s.server.lastRequest.OrderUuid)
	s.Equal("user-1", s.server.lastRequest.UserUuid)
	s.Equal(paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, s.server.lastRequest.PaymentMethod)
}

func (s *PaymentClientSuite) TestPayOrder_InvalidPaymentMethod() {
	transactionUUID, err := s.client.PayOrder(s.ctx, "order-1", "user-1", model.PaymentMethodUnspecified)

	s.Require().Error(err)
	s.Empty(transactionUUID)
	s.Nil(s.server.lastRequest)
}

func (s *PaymentClientSuite) TestPayOrder_ServerError() {
	s.server.payOrderErr = status.Error(codes.Internal, "payment unavailable")

	transactionUUID, err := s.client.PayOrder(s.ctx, "order-1", "user-1", model.PaymentMethodCard)

	s.Require().Error(err)
	s.Empty(transactionUUID)
}

func TestPaymentClientSuite(t *testing.T) {
	suite.Run(t, new(PaymentClientSuite))
}

// paymentServiceServer — тестовая реализация paymentv1.PaymentServiceServer
type paymentServiceServer struct {
	paymentv1.UnimplementedPaymentServiceServer
	payOrderResponse *paymentv1.PayOrderResponse
	payOrderErr      error
	lastRequest      *paymentv1.PayOrderRequest
}

func newPaymentServiceServer() *paymentServiceServer {
	return &paymentServiceServer{}
}

func (s *paymentServiceServer) PayOrder(_ context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	s.lastRequest = req

	if s.payOrderErr != nil {
		return nil, s.payOrderErr
	}

	return s.payOrderResponse, nil
}
