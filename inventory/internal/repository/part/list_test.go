package part_test

import (
	"context"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *PartRepositorySuite) TestListParts_ReturnsAllParts() {
	first := testPart()
	second := testPart()
	second.UUID = "uuid-2"
	second.Name = "Turbo Charger"

	s.collection.
		EXPECT().
		Find(s.ctx, mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, _, results any) error {
			*results.(*[]repositoryModel.Part) = []repositoryModel.Part{first, second}
			return nil
		}).
		Once()

	parts, err := s.repo.ListParts(s.ctx, model.PartsFilter{})

	require.NoError(s.T(), err)
	require.Len(s.T(), parts, 2)
	assert.Equal(s.T(), first.UUID, parts[0].UUID)
	assert.Equal(s.T(), second.UUID, parts[1].UUID)
}

func (s *PartRepositorySuite) TestListParts_FiltersByUUID() {
	want := testPart()

	s.collection.
		EXPECT().
		Find(s.ctx, mock.MatchedBy(func(filter any) bool {
			got, ok := filter.(bson.M)
			if !ok {
				return false
			}
			in, ok := got["uuid"].(bson.M)
			if !ok {
				return false
			}
			return assert.ElementsMatch(s.T(), []string{want.UUID}, in["$in"].([]string))
		}), mock.Anything).
		RunAndReturn(func(_ context.Context, _, results any) error {
			*results.(*[]repositoryModel.Part) = []repositoryModel.Part{want}
			return nil
		}).
		Once()

	filtered, err := s.repo.ListParts(s.ctx, model.PartsFilter{UUIDs: []string{want.UUID}})

	require.NoError(s.T(), err)
	require.Len(s.T(), filtered, 1)
	require.Equal(s.T(), want.UUID, filtered[0].UUID)
}

func (s *PartRepositorySuite) TestListParts_FiltersByNames() {
	want := testPart()

	s.collection.
		EXPECT().
		Find(s.ctx, mock.MatchedBy(func(filter any) bool {
			got, ok := filter.(bson.M)
			if !ok {
				return false
			}
			in, ok := got["name"].(bson.M)
			if !ok {
				return false
			}
			return assert.ElementsMatch(s.T(), []string{want.Name}, in["$in"].([]string))
		}), mock.Anything).
		RunAndReturn(func(_ context.Context, _, results any) error {
			*results.(*[]repositoryModel.Part) = []repositoryModel.Part{want}
			return nil
		}).
		Once()

	filtered, err := s.repo.ListParts(s.ctx, model.PartsFilter{Names: []string{want.Name}})

	require.NoError(s.T(), err)
	require.Len(s.T(), filtered, 1)
	require.Equal(s.T(), want.Name, filtered[0].Name)
}

func (s *PartRepositorySuite) TestListParts_FiltersByTags() {
	want := testPart()

	s.collection.
		EXPECT().
		Find(s.ctx, mock.MatchedBy(func(filter any) bool {
			got, ok := filter.(bson.M)
			if !ok {
				return false
			}
			in, ok := got["tags"].(bson.M)
			if !ok {
				return false
			}
			return assert.ElementsMatch(s.T(), []string{"high-performance"}, in["$in"].([]string))
		}), mock.Anything).
		RunAndReturn(func(_ context.Context, _, results any) error {
			*results.(*[]repositoryModel.Part) = []repositoryModel.Part{want}
			return nil
		}).
		Once()

	filtered, err := s.repo.ListParts(s.ctx, model.PartsFilter{Tags: []string{"high-performance"}})

	require.NoError(s.T(), err)
	require.Len(s.T(), filtered, 1)
}

func (s *PartRepositorySuite) TestListParts_EmptyFilterQuery() {
	first := testPart()

	s.collection.
		EXPECT().
		Find(s.ctx, mock.MatchedBy(func(filter any) bool {
			got, ok := filter.(bson.M)
			if !ok {
				return false
			}
			return len(got) == 0
		}), mock.Anything).
		RunAndReturn(func(_ context.Context, _, results any) error {
			*results.(*[]repositoryModel.Part) = []repositoryModel.Part{first}
			return nil
		}).
		Once()

	parts, err := s.repo.ListParts(s.ctx, model.PartsFilter{})

	require.NoError(s.T(), err)
	require.Len(s.T(), parts, 1)
}
