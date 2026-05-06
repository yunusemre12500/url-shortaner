package factory

import (
	"github.com/spf13/pflag"
	"go.yunus-emre.dev/url-shortaner/storage/postgres"
)

const (
	BackendPostgres = "postgres"
)

type Config struct {
	Backend  string
	Postgres *postgres.Config

	// prevents unkeyed literal initialization
	_ struct{}
}

var DefaultConfig = &Config{
	Backend:  BackendPostgres,
	Postgres: postgres.DefaultConfig,
}

func (config *Config) ApplyFlagsWithPrefix(prefix string, flagSet *pflag.FlagSet) {
	flagSet.StringVar(&config.Backend, prefix+"backend", DefaultConfig.Backend, "")
	config.Postgres.ApplyFlagsWithPrefix(prefix+"postgres.", flagSet)
}
