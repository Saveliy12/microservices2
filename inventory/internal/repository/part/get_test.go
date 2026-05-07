package part_test

import (
	"context"
	"testing"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repoinventory "github.com.Saveliy12/microservices2.git/inventory/internal/repository/part"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PartRepositorySuite struct {
	suite.Suite
	repo *repoinventory.Repository
}

func (s *PartRepositorySuite) SetupTest() {
	s.repo = repoinventory.NewRepository()
}
func (s *PartRepositorySuite) TestGetPart_Success() {
	s.T().Parallel()

	parts, err := s.repo.ListParts(context.Background(), model.PartsFilter{})
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), parts)

	want := parts[0]
	got, err := s.repo.GetPart(context.Background(), want.UUID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), want.UUID, got.UUID)
	assert.Equal(s.T(), want.Name, got.Name)
	assert.Equal(s.T(), want.Price, got.Price)
}

func (s *PartRepositorySuite) TestGetPart_NotFound() {
	s.T().Parallel()

	_, err := s.repo.GetPart(context.Background(), "missing-uuid")

	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrPartNotFound)
}

func TestPartRepositorySuite(t *testing.T) {
	suite.Run(t, new(PartRepositorySuite))
}
