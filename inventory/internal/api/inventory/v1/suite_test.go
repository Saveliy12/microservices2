package v1

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestInventoryGetAPIIntegration(t *testing.T) {
	suite.Run(t, new(InventoryGetAPISuite))
}

func TestInventoryListAPIIntegration(t *testing.T) {
	suite.Run(t, new(InventoryListAPISuite))
}
