package payment_test

import (
	"context"
	"regexp"
	"testing"

	"github.com.Saveliy12/microservices2.git/payment/internal/model"
	paymentservice "github.com.Saveliy12/microservices2.git/payment/internal/service/payment"
	"github.com/stretchr/testify/suite"
)

type PayOrderServiceTestSuite struct {
	suite.Suite
	service *paymentservice.Service
}

func (s *PayOrderServiceTestSuite) SetupTest() {
	// Инициализируем сервис перед каждым тестом
	s.service = paymentservice.NewService()
}

func (s *PayOrderServiceTestSuite) TestPayOrder_Success() {
	// Arrange
	ctx := context.Background()
	payment := model.Payment{
		OrderUUID:     "order-123",
		UserUUID:      "user-123",
		PaymentMethod: "CARD",
	}

	txnUUID, err := s.service.PayOrder(ctx, payment)

	s.NoError(err)
	s.NotEmpty(txnUUID)
}

func (s *PayOrderServiceTestSuite) TestPayOrder_ReturnsValidUUID() {
	ctx := context.Background()
	payment := model.Payment{
		OrderUUID:     "order-456",
		UserUUID:      "user-456",
		PaymentMethod: "SBP",
	}

	txnUUID, err := s.service.PayOrder(ctx, payment)

	s.NoError(err)

	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	s.True(uuidRegex.MatchString(txnUUID), "Returned string is not a valid UUID: %s", txnUUID)
}

func (s *PayOrderServiceTestSuite) TestPayOrder_GeneratesUniqueTransactionUUIDs() {
	ctx := context.Background()
	payment := model.Payment{
		OrderUUID:     "order-789",
		UserUUID:      "user-789",
		PaymentMethod: "CREDIT_CARD",
	}

	txnUUID1, err1 := s.service.PayOrder(ctx, payment)
	txnUUID2, err2 := s.service.PayOrder(ctx, payment)
	txnUUID3, err3 := s.service.PayOrder(ctx, payment)

	s.NoError(err1)
	s.NoError(err2)
	s.NoError(err3)

	s.NotEqual(txnUUID1, txnUUID2, "Transaction UUIDs should be unique")
	s.NotEqual(txnUUID2, txnUUID3, "Transaction UUIDs should be unique")
	s.NotEqual(txnUUID1, txnUUID3, "Transaction UUIDs should be unique")
}

func (s *PayOrderServiceTestSuite) TestPayOrder_WorksWithDifferentPaymentMethods() {
	ctx := context.Background()
	paymentMethods := []string{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}

	for _, method := range paymentMethods {
		payment := model.Payment{
			OrderUUID:     "order-" + method,
			UserUUID:      "user-" + method,
			PaymentMethod: method,
		}

		txnUUID, err := s.service.PayOrder(ctx, payment)

		s.NoError(err, "Payment method %s should succeed", method)
		s.NotEmpty(txnUUID, "Transaction UUID should not be empty for %s", method)
	}
}

func (s *PayOrderServiceTestSuite) TestPayOrder_ProcessesPaymentWithMinimalData() {
	ctx := context.Background()
	payment := model.Payment{
		OrderUUID:     "order-min",
		UserUUID:      "user-min",
		PaymentMethod: "CARD",
	}

	txnUUID, err := s.service.PayOrder(ctx, payment)

	s.NoError(err)
	s.NotEmpty(txnUUID)
}

func (s *PayOrderServiceTestSuite) TestPayOrder_HandlesContextCancellation() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // отменяем контекст

	payment := model.Payment{
		OrderUUID:     "order-ctx",
		UserUUID:      "user-ctx",
		PaymentMethod: "CARD",
	}

	txnUUID, err := s.service.PayOrder(ctx, payment)

	s.NoError(err)
	s.NotEmpty(txnUUID)
}

func (s *PayOrderServiceTestSuite) TestPayOrder_ConcurrentPayments() {
	ctx := context.Background()
	txnUUIDs := make(chan string, 3)

	for i := range 3 {
		go func(index int) {
			payment := model.Payment{
				OrderUUID:     "order-concurrent-" + string(rune(index)),
				UserUUID:      "user-concurrent-" + string(rune(index)),
				PaymentMethod: "CARD",
			}
			txnUUID, _ := s.service.PayOrder(ctx, payment)
			txnUUIDs <- txnUUID
		}(i)
	}

	uuids := make(map[string]bool)
	for range 3 {
		uuid := <-txnUUIDs
		s.NotEmpty(uuid)
		s.False(uuids[uuid], "Transaction UUIDs should be unique in concurrent calls")
		uuids[uuid] = true
	}

	s.Equal(3, len(uuids), "Should have 3 unique UUIDs")
}

func TestPayOrderServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PayOrderServiceTestSuite))
}
