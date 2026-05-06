package tracing

import "github.com/spf13/pflag"

type Config struct {
	Batching    *BatchingConfig
	Endpoint    string
	Insecure    bool
	Protocol    string
	ServiceName string

	// prevents unkeyed literal initialization
	_ struct{}
}

var DefaultConfig = &Config{
	Batching:    DefaultBatchingConfig,
	Endpoint:    "localhost:4317",
	Insecure:    false,
	Protocol:    "grpc",
	ServiceName: "api-server",
}

func (config *Config) ApplyFlagsWithPrefix(prefix string, flagSet *pflag.FlagSet) {
	config.Batching.ApplyFlagsWithPrefix(prefix+"batching.", flagSet)
	flagSet.StringVar(&config.Endpoint, prefix+"endpoint", DefaultConfig.Endpoint, "")
	flagSet.BoolVar(&config.Insecure, prefix+"insecure", DefaultConfig.Insecure, "")
	flagSet.StringVar(&config.Protocol, prefix+"protocol", DefaultConfig.Protocol, "")
}

type BatchingConfig struct {
	Enabled bool

	// prevents unkeyed literal initialization
	_ struct{}
}

var DefaultBatchingConfig = &BatchingConfig{
	Enabled: false,
}

func (config *BatchingConfig) ApplyFlagsWithPrefix(prefix string, flagSet *pflag.FlagSet) {
	flagSet.BoolVar(&config.Enabled, prefix+"enabled", DefaultBatchingConfig.Enabled, "")
}
