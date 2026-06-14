package part_test

import (
	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	"github.com/stretchr/testify/require"
)

func (s *PartServiceSuite) TestListParts_Success() {
	filter := model.PartsFilter{Names: []string{"Engine V8"}}
	want := []model.Part{
		{UUID: "uuid-1", Name: "Engine V8", Price: 50000},
	}

	s.partRepositoryMock.EXPECT().
		ListParts(s.ctx, filter).
		Return(want, nil).
		Once()

	got, err := s.service.ListParts(s.ctx, filter)

	require.NoError(s.T(), err)
	require.Equal(s.T(), want, got)
}

func (s *PartServiceSuite) TestListParts_ReturnsError() {
	filter := model.PartsFilter{UUIDs: []string{"missing-uuid"}}

	s.partRepositoryMock.EXPECT().
		ListParts(s.ctx, filter).
		Return(nil, model.ErrPartNotFound).
		Once()

	got, err := s.service.ListParts(s.ctx, filter)

	require.ErrorIs(s.T(), err, model.ErrPartNotFound)
	require.Nil(s.T(), got)
}
