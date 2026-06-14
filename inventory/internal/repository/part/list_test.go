package part_test

import (
	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	"github.com/stretchr/testify/require"
)

func (s *PartRepositorySuite) TestListParts_ReturnsAllSeedParts() {
	parts, err := s.repo.ListParts(s.ctx, model.PartsFilter{})

	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), parts)

}

func (s *PartRepositorySuite) TestListParts_FiltersByUUID() {
	parts, err := s.repo.ListParts(s.ctx, model.PartsFilter{})
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), parts)

	want := parts[0]
	filtered, err := s.repo.ListParts(s.ctx, model.PartsFilter{UUIDs: []string{want.UUID}})

	require.NoError(s.T(), err)
	require.Len(s.T(), filtered, 1)
	require.Equal(s.T(), want.UUID, filtered[0].UUID)
}
