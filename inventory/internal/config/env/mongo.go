package env

import (
	"errors"
	"os"
)

const (
	mongoURIEnvName = "MONGO_URI"
	mongoDBEnvName  = "MONGO_DB"
)

type mongoConfig struct {
	uri      string
	database string
}

func NewMongoConfig() (*mongoConfig, error) {
	uri := os.Getenv(mongoURIEnvName)
	if uri == "" {
		return nil, errors.New("MONGO_URI is not set")
	}
	db := os.Getenv(mongoDBEnvName)
	if db == "" {
		return nil, errors.New("MONGO_DB is not set")
	}
	return &mongoConfig{uri: uri, database: db}, nil
}

func (c *mongoConfig) URI() string {
	return c.uri
}

func (c *mongoConfig) Database() string {
	return c.database
}
