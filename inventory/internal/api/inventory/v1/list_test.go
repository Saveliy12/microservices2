package v1

import (
	"context"
	"errors"
	"testing"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	serviceMocks "github.com.Saveliy12/microservices2.git/inventory/internal/service/mocks"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type InventoryListAPISuite struct {
	suite.Suite
	api           *API
	inventoryMock *serviceMocks.InventoryService
}

func (s *InventoryListAPISuite) SetupTest() {
	s.inventoryMock = serviceMocks.NewInventoryService(s.T())
	s.api = NewAPI(s.inventoryMock)
}

func (s *InventoryListAPISuite) TestListParts_NilRequest() {
	resp, err := s.api.ListParts(context.Background(), nil)

	require.Nil(s.T(), resp)
	require.Error(s.T(), err)
	st, _ := status.FromError(err)
	assert.Equal(s.T(), codes.InvalidArgument, st.Code())
}

func (s *InventoryListAPISuite) TestListParts_Success() {
	filter := &inventoryv1.PartsFilter{
		Names: []*wrapperspb.StringValue{{Value: "Engine V8"}},
	}
	want := []model.Part{{UUID: "uuid-1", Name: "Engine V8", Price: 50000}}

	s.inventoryMock.EXPECT().
		ListParts(context.Background(), model.PartsFilter{
			UUIDs:                 []string{},
			Names:                 []string{"Engine V8"},
			Categories:            []model.Category{},
			ManufacturerCountries: []string{},
			Tags:                  []string{},
		}).
		Return(want, nil).
		Once()

	resp, err := s.api.ListParts(context.Background(), &inventoryv1.ListPartsRequest{Filter: filter})

	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
	require.Len(s.T(), resp.Parts, 1)
	assert.Equal(s.T(), want[0].UUID, resp.Parts[0].Uuid)
	assert.Equal(s.T(), want[0].Name, resp.Parts[0].Name)
	assert.Equal(s.T(), want[0].Price, resp.Parts[0].Price)
}

func (s *InventoryListAPISuite) TestListParts_ServiceError() {
	wantErr := errors.New("service failed")

	s.inventoryMock.EXPECT().
		ListParts(context.Background(), model.PartsFilter{}).
		Return(nil, wantErr).
		Once()

	resp, err := s.api.ListParts(context.Background(), &inventoryv1.ListPartsRequest{})

	require.Nil(s.T(), resp)
	require.Error(s.T(), err)
	st, _ := status.FromError(err)
	assert.Equal(s.T(), codes.Internal, st.Code())
	assert.Equal(s.T(), wantErr.Error(), st.Message())
}

func TestInventoryListAPISuite(t *testing.T) {
	suite.Run(t, new(InventoryListAPISuite))
}
