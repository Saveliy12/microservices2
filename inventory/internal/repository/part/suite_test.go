package part_test

import (
	"context"
	"testing"

	dbbmocks "github.com.Saveliy12/microservices2.git/inventory/internal/client/db/mocks"
	repoinventory "github.com.Saveliy12/microservices2.git/inventory/internal/repository/part"
	"github.com/stretchr/testify/suite"
)

type PartRepositorySuite struct {
	suite.Suite
	ctx        context.Context
	collection *dbbmocks.CollectionMock
	repo       *repoinventory.Repository
}

func (s *PartRepositorySuite) SetupTest() {
	s.ctx = context.Background()
	s.collection = dbbmocks.NewCollectionMock(s.T())
	s.repo = repoinventory.NewRepository(s.collection)
}

func (s *PartRepositorySuite) TearDownTest() {
	// Cleanup if needed
}

func TestPartRepositorySuite(t *testing.T) {
	suite.Run(t, new(PartRepositorySuite))
}
