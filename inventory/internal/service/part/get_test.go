package part_test

import (
	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	"github.com/stretchr/testify/require"
)

func (s *PartServiceSuite) TestGetPart_Success() {
	uuid := "550e8400-e29b-41d4-a716-446655440000"
	want := model.Part{UUID: uuid, Name: "Engine V8", Price: 50000}

	s.partRepositoryMock.EXPECT().
		GetPart(s.ctx, uuid).
		Return(want, nil).
		Once()

	got, err := s.service.GetPart(s.ctx, uuid)

	require.NoError(s.T(), err)
	require.Equal(s.T(), want, got)
}

func (s *PartServiceSuite) TestGetPart_ReturnsError() {
	uuid := "missing-uuid"
	wantErr := model.ErrPartNotFound

	s.partRepositoryMock.EXPECT().
		GetPart(s.ctx, uuid).
		Return(model.Part{}, wantErr).
		Once()

	got, err := s.service.GetPart(s.ctx, uuid)

	require.ErrorIs(s.T(), err, wantErr)
	require.Equal(s.T(), model.Part{}, got)
}
