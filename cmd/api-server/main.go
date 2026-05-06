package main

import (
	"os"

	"github.com/spf13/pflag"
	"go.yunus-emre.dev/url-shortaner/apiserver"
)

func main() {
	flagSet := pflag.NewFlagSet("api-server", pflag.ExitOnError)

	config := apiserver.DefaultConfig

	config.ApplyFlags(flagSet)

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	apiServer, err := apiserver.New(config)

	if err != nil {
		panic(err)
	}

	apiServer.Run()
}
