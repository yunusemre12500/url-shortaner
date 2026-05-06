package apiserver

import (
	"time"

	"github.com/spf13/pflag"
	"go.yunus-emre.dev/url-shortaner/internal/tracing"
	storagefactory "go.yunus-emre.dev/url-shortaner/storage/factory"
)

type Config struct {
	GracefulShutdownTimeout time.Duration
	IdleTimeout             time.Duration
	MaxHeaderBytes          int
	ReadHeaderTimeout       time.Duration
	ReadTimeout             time.Duration
	Storage                 *storagefactory.Config
	ListenAddress           string
	Tracing                 *tracing.Config
	WriteTimeout            time.Duration

	// prevents unkeyed literal initialization
	_ struct{}
}

var DefaultConfig = &Config{
	GracefulShutdownTimeout: 60 * time.Second,
	IdleTimeout:             60 * time.Second,
	MaxHeaderBytes:          4096,
	ReadHeaderTimeout:       10 * time.Second,
	ReadTimeout:             30 * time.Second,
	Storage:                 storagefactory.DefaultConfig,
	ListenAddress:           "0.0.0.0:80",
	Tracing:                 tracing.DefaultConfig,
	WriteTimeout:            30 * time.Second,
}

func (config *Config) ApplyFlags(flagSet *pflag.FlagSet) {
	flagSet.DurationVar(&config.GracefulShutdownTimeout, "graceful-shutdown-timeout", DefaultConfig.GracefulShutdownTimeout, "")
	flagSet.DurationVar(&config.IdleTimeout, "idle-timeout", DefaultConfig.IdleTimeout, "")
	flagSet.IntVar(&config.MaxHeaderBytes, "max-header-bytes", DefaultConfig.MaxHeaderBytes, "")
	flagSet.DurationVar(&config.ReadHeaderTimeout, "read-header-timeout", DefaultConfig.ReadHeaderTimeout, "")
	flagSet.DurationVar(&config.ReadTimeout, "read-timeout", DefaultConfig.ReadTimeout, "")
	config.Storage.ApplyFlagsWithPrefix("storage.", flagSet)
	flagSet.StringVar(&config.ListenAddress, "listen-address", "0.0.0.0:80", "")
	config.Tracing.ApplyFlagsWithPrefix("tracing.", flagSet)
	flagSet.DurationVar(&config.WriteTimeout, "write-timeout", DefaultConfig.WriteTimeout, "")
}
