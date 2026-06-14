package v1

import (
	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *InventoryAPISuite) TestGetPart_NilRequest() {
	resp, err := s.api.GetPart(s.ctx, nil)

	require.Nil(s.T(), resp)
	require.Error(s.T(), err)
	st, _ := status.FromError(err)
	assert.Equal(s.T(), codes.InvalidArgument, st.Code())
}

func (s *InventoryAPISuite) TestGetPart_Success() {
	uuid := "550e8400-e29b-41d4-a716-446655440000"
	want := model.Part{UUID: uuid, Name: "Engine V8", Price: 50000}

	s.inventoryMock.EXPECT().
		GetPart(s.ctx, uuid).
		Return(want, nil).
		Once()

	resp, err := s.api.GetPart(s.ctx, &inventoryv1.GetPartRequest{Uuid: uuid})

	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
	require.NotNil(s.T(), resp.Part)
	assert.Equal(s.T(), uuid, resp.Part.Uuid)
	assert.Equal(s.T(), want.Name, resp.Part.Name)
	assert.Equal(s.T(), want.Price, resp.Part.Price)
}

func (s *InventoryAPISuite) TestGetPart_NotFound() {
	uuid := "missing-uuid"
	wantErr := model.ErrPartNotFound

	s.inventoryMock.EXPECT().
		GetPart(s.ctx, uuid).
		Return(model.Part{}, wantErr).
		Once()

	resp, err := s.api.GetPart(s.ctx, &inventoryv1.GetPartRequest{Uuid: uuid})

	require.Nil(s.T(), resp)
	require.Error(s.T(), err)
	st, _ := status.FromError(err)
	assert.Equal(s.T(), codes.NotFound, st.Code())
	assert.Equal(s.T(), wantErr.Error(), st.Message())
}
