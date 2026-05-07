package part_test

import (
	"context"
	"errors"
	"testing"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repositoryMocks "github.com.Saveliy12/microservices2.git/inventory/internal/repository/mocks"
	servicepart "github.com.Saveliy12/microservices2.git/inventory/internal/service/part"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PartServiceSuite struct {
	suite.Suite
	service            *servicepart.Service
	partRepositoryMock *repositoryMocks.PartRepository
}

func (s *PartServiceSuite) SetupTest() {
	s.partRepositoryMock = repositoryMocks.NewPartRepository(s.T())
	s.service = servicepart.NewService(s.partRepositoryMock)
}

func (s *PartServiceSuite) TestGetPart_Success() {
	uuid := "550e8400-e29b-41d4-a716-446655440000"
	want := model.Part{UUID: uuid, Name: "Engine V8", Price: 50000}

	s.partRepositoryMock.EXPECT().
		GetPart(context.Background(), uuid).
		Return(want, nil).
		Once()

	got, err := s.service.GetPart(context.Background(), uuid)

	require.NoError(s.T(), err)
	require.Equal(s.T(), want, got)
}

func (s *PartServiceSuite) TestGetPart_ReturnsError() {
	uuid := "missing-uuid"
	wantErr := errors.New("part not found")

	s.partRepositoryMock.EXPECT().
		GetPart(context.Background(), uuid).
		Return(model.Part{}, wantErr).
		Once()

	got, err := s.service.GetPart(context.Background(), uuid)

	require.ErrorIs(s.T(), err, wantErr)
	require.Equal(s.T(), model.Part{}, got)
}

func TestPartServiceSuite(t *testing.T) {
	suite.Run(t, new(PartServiceSuite))
}
