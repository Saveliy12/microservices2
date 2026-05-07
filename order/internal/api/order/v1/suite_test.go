package v1

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestOrderGetAPIIntegration(t *testing.T) {
	suite.Run(t, new(GetOrderAPISuite))
}

func TestOrderCreateAPIIntegration(t *testing.T) {
	suite.Run(t, new(CreateOrderAPISuite))
}

func TestOrderPayAPIIntegration(t *testing.T) {
	suite.Run(t, new(PayOrderAPISuite))
}

func TestOrderCancelAPIIntegration(t *testing.T) {
	suite.Run(t, new(CancelOrderAPISuite))
}
