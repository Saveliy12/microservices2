package part_test

import (
	"context"
	"testing"

	repositoryMocks "github.com.Saveliy12/microservices2.git/inventory/internal/repository/mocks"
	servicepart "github.com.Saveliy12/microservices2.git/inventory/internal/service/part"
	"github.com/stretchr/testify/suite"
)

type PartServiceSuite struct {
	suite.Suite
	ctx                context.Context
	service            *servicepart.Service
	partRepositoryMock *repositoryMocks.PartRepositoryMock
}

func (s *PartServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.partRepositoryMock = repositoryMocks.NewPartRepositoryMock(s.T())
	s.service = servicepart.NewService(s.partRepositoryMock)
}

func TestPartServiceSuite(t *testing.T) {
	suite.Run(t, new(PartServiceSuite))
}
