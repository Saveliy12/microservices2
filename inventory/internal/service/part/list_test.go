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

type PartListServiceSuite struct {
	suite.Suite
	service            *servicepart.Service
	partRepositoryMock *repositoryMocks.PartRepository
}

func (s *PartListServiceSuite) SetupTest() {
	s.partRepositoryMock = repositoryMocks.NewPartRepository(s.T())
	s.service = servicepart.NewService(s.partRepositoryMock)
}

func (s *PartListServiceSuite) TestListParts_Success() {
	filter := model.PartsFilter{Names: []string{"Engine V8"}}
	want := []model.Part{
		{UUID: "uuid-1", Name: "Engine V8", Price: 50000},
	}

	s.partRepositoryMock.EXPECT().
		ListParts(context.Background(), filter).
		Return(want, nil).
		Once()

	got, err := s.service.ListParts(context.Background(), filter)

	require.NoError(s.T(), err)
	require.Equal(s.T(), want, got)
}

func (s *PartListServiceSuite) TestListParts_ReturnsError() {
	filter := model.PartsFilter{UUIDs: []string{"missing-uuid"}}
	wantErr := errors.New("repository failed")

	s.partRepositoryMock.EXPECT().
		ListParts(context.Background(), filter).
		Return(nil, wantErr).
		Once()

	got, err := s.service.ListParts(context.Background(), filter)

	require.ErrorIs(s.T(), err, wantErr)
	require.Nil(s.T(), got)
}

func TestPartListServiceSuite(t *testing.T) {
	suite.Run(t, new(PartListServiceSuite))
}
