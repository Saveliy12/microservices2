package part_test

import (
	"context"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *PartRepositorySuite) TestGetPart_Success() {
	want := testPart()
	s.collection.
		EXPECT().
		FindOne(s.ctx, mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, _, result any) error {
			*result.(*repositoryModel.Part) = want
			return nil
		}).
		Once()

	got, err := s.repo.GetPart(s.ctx, want.UUID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), want.UUID, got.UUID)
	assert.Equal(s.T(), want.Name, got.Name)
	assert.Equal(s.T(), want.Price, got.Price)
}

func (s *PartRepositorySuite) TestGetPart_NotFound() {
	s.collection.
		EXPECT().
		FindOne(s.ctx, mock.Anything, mock.Anything).
		Return(mongo.ErrNoDocuments).
		Once()

	_, err := s.repo.GetPart(s.ctx, "missing-uuid")

	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrPartNotFound)
}

func (s *PartRepositorySuite) TestGetPart_DBError() {
	s.collection.
		EXPECT().
		FindOne(s.ctx, mock.Anything, mock.Anything).
		Return(mongo.ErrClientDisconnected).
		Once()

	_, err := s.repo.GetPart(s.ctx, "uuid-1")

	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, mongo.ErrClientDisconnected)
}
