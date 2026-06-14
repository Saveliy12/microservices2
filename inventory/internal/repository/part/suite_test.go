package part_test

import (
	"context"
	"testing"

	repoinventory "github.com.Saveliy12/microservices2.git/inventory/internal/repository/part"
	"github.com/stretchr/testify/suite"
)

type PartRepositorySuite struct {
	suite.Suite
	ctx  context.Context
	repo *repoinventory.Repository
}

func (s *PartRepositorySuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = repoinventory.NewRepository()
}

func (s *PartRepositorySuite) TearDownTest() {
	// Cleanup if needed
}

func TestPartRepositorySuite(t *testing.T) {
	suite.Run(t, new(PartRepositorySuite))
}
