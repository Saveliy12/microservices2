package v1

import (
	"context"
	"testing"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repositoryMocks "github.com.Saveliy12/microservices2.git/inventory/internal/service/mocks"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryGetAPISuite struct {
	suite.Suite
	api           *API
	inventoryMock *repositoryMocks.InventoryService
}

func (s *InventoryGetAPISuite) SetupTest() {
	s.inventoryMock = repositoryMocks.NewInventoryService(s.T())
	s.api = NewAPI(s.inventoryMock)
}

func (s *InventoryGetAPISuite) TestGetPart_NilRequest() {
	resp, err := s.api.GetPart(context.Background(), nil)

	require.Nil(s.T(), resp)
	require.Error(s.T(), err)
	st, _ := status.FromError(err)
	assert.Equal(s.T(), codes.InvalidArgument, st.Code())
}

func (s *InventoryGetAPISuite) TestGetPart_Success() {
	uuid := "550e8400-e29b-41d4-a716-446655440000"
	want := model.Part{UUID: uuid, Name: "Engine V8", Price: 50000}

	s.inventoryMock.EXPECT().
		GetPart(context.Background(), uuid).
		Return(want, nil).
		Once()

	resp, err := s.api.GetPart(context.Background(), &inventoryv1.GetPartRequest{Uuid: uuid})

	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
	require.NotNil(s.T(), resp.Part)
	assert.Equal(s.T(), uuid, resp.Part.Uuid)
	assert.Equal(s.T(), want.Name, resp.Part.Name)
	assert.Equal(s.T(), want.Price, resp.Part.Price)
}

func (s *InventoryGetAPISuite) TestGetPart_NotFound() {
	uuid := "missing-uuid"
	wantErr := model.ErrPartNotFound

	s.inventoryMock.EXPECT().
		GetPart(context.Background(), uuid).
		Return(model.Part{}, wantErr).
		Once()

	resp, err := s.api.GetPart(context.Background(), &inventoryv1.GetPartRequest{Uuid: uuid})

	require.Nil(s.T(), resp)
	require.Error(s.T(), err)
	st, _ := status.FromError(err)
	assert.Equal(s.T(), codes.NotFound, st.Code())
	assert.Equal(s.T(), wantErr.Error(), st.Message())
}

func TestInventoryGetAPISuite(t *testing.T) {
	suite.Run(t, new(InventoryGetAPISuite))
}
