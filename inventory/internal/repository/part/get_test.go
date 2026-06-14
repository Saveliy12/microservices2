package part_test

import (
	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PartRepositorySuite) TestGetPart_Success() {

	parts, err := s.repo.ListParts(s.ctx, model.PartsFilter{})
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), parts)

	want := parts[0]
	got, err := s.repo.GetPart(s.ctx, want.UUID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), want.UUID, got.UUID)
	assert.Equal(s.T(), want.Name, got.Name)
	assert.Equal(s.T(), want.Price, got.Price)
}

func (s *PartRepositorySuite) TestGetPart_NotFound() {

	_, err := s.repo.GetPart(s.ctx, "missing-uuid")

	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrPartNotFound)
}
