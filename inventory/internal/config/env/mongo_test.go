package env_test

import (
	"testing"

	env "github.com.Saveliy12/microservices2.git/inventory/internal/config/env"
	"github.com/stretchr/testify/suite"
)

type MongoConfigTestSuite struct {
	suite.Suite
}

func (s *MongoConfigTestSuite) TestNewMongoConfig_Default() {
	t := s.T()
	t.Setenv("MONGO_URI", "mongodb://localhost:27017")
	t.Setenv("MONGO_DB", "testdb")

	cfg, err := env.NewMongoConfig()

	s.Require().NoError(err)
	s.Require().Equal("mongodb://localhost:27017", cfg.URI())
	s.Require().Equal("testdb", cfg.Database())
}

func (s *MongoConfigTestSuite) TestNewMongoConfig_MissingURI() {
	t := s.T()
	t.Setenv("MONGO_URI", "")
	t.Setenv("MONGO_DB", "testdb")

	cfg, err := env.NewMongoConfig()

	s.Require().Error(err)
	s.Nil(cfg)
}

func (s *MongoConfigTestSuite) TestNewMongoConfig_MissingDB() {
	t := s.T()
	t.Setenv("MONGO_URI", "mongodb://localhost:27017")
	t.Setenv("MONGO_DB", "")

	cfg, err := env.NewMongoConfig()

	s.Require().Error(err)
	s.Nil(cfg)
}

func TestMongoConfigTestSuite(t *testing.T) {
	suite.Run(t, new(MongoConfigTestSuite))
}
