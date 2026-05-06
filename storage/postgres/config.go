package postgres

import (
	"time"

	"github.com/spf13/pflag"
)

type Config struct {
	ConnectionIdleTimeout time.Duration
	ConnectionMaxLifetime time.Duration
	DSN                   string
	MaxIdleConnections    int
	MaxOpenConnections    int

	// prevents unkeyed literal initialization
	_ struct{}
}

var DefaultConfig = &Config{
	ConnectionIdleTimeout: 60 * time.Second,
	ConnectionMaxLifetime: 1 * time.Hour,
	DSN:                   "postgres://localhost:5432",
	MaxIdleConnections:    3,
	MaxOpenConnections:    10,
}

func (config *Config) ApplyFlagsWithPrefix(prefix string, flagSet *pflag.FlagSet) {
	flagSet.DurationVar(&config.ConnectionIdleTimeout, prefix+"connection-idle-timeout", DefaultConfig.ConnectionIdleTimeout, "")
	flagSet.DurationVar(&config.ConnectionIdleTimeout, prefix+"connection-max-lifetime", DefaultConfig.ConnectionMaxLifetime, "")
	flagSet.StringVar(&config.DSN, prefix+"dsn", DefaultConfig.DSN, "")
	flagSet.IntVar(&config.MaxIdleConnections, prefix+"max-idle-connections", DefaultConfig.MaxIdleConnections, "")
	flagSet.IntVar(&config.MaxOpenConnections, prefix+"max-open-connections", DefaultConfig.MaxOpenConnections, "")
}
