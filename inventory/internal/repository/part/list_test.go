package part_test

import (
	"context"
	"testing"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repoinventory "github.com.Saveliy12/microservices2.git/inventory/internal/repository/part"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PartRepositoryListSuite struct {
	suite.Suite
	repo *repoinventory.Repository
}

func (s *PartRepositoryListSuite) SetupTest() {
	s.repo = repoinventory.NewRepository()
}

func (s *PartRepositoryListSuite) TestListParts_ReturnsAllSeedParts() {
	parts, err := s.repo.ListParts(context.Background(), model.PartsFilter{})

	require.NoError(s.T(), err)
	require.Len(s.T(), parts, 4)
}

func (s *PartRepositoryListSuite) TestListParts_FiltersByUUID() {
	parts, err := s.repo.ListParts(context.Background(), model.PartsFilter{})
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), parts)

	want := parts[0]
	filtered, err := s.repo.ListParts(context.Background(), model.PartsFilter{UUIDs: []string{want.UUID}})

	require.NoError(s.T(), err)
	require.Len(s.T(), filtered, 1)
	require.Equal(s.T(), want.UUID, filtered[0].UUID)
}

func TestPartRepositoryListSuite(t *testing.T) {
	suite.Run(t, new(PartRepositoryListSuite))
}
