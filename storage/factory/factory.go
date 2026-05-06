package factory

import (
	"errors"
	"fmt"

	"go.yunus-emre.dev/url-shortaner/storage"
	"go.yunus-emre.dev/url-shortaner/storage/postgres"
)

func Create(config *Config) (storage.Storage, error) {
	switch config.Backend {
	case "":
		return nil, errors.New("storage backend cannot be empty")
	case BackendPostgres:
		return postgres.New(config.Postgres)
	default:
		return nil, fmt.Errorf("storage backend '%s' not supported.", config.Backend)
	}
}
