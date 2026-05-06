package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.yunus-emre.dev/url-shortaner/apiserver/internal/api/rest"
	"go.yunus-emre.dev/url-shortaner/internal/tracing"
	"go.yunus-emre.dev/url-shortaner/pkg/version"
	"go.yunus-emre.dev/url-shortaner/storage"
	storagefactory "go.yunus-emre.dev/url-shortaner/storage/factory"
)

type APIServer struct {
	config     *Config
	httpServer *http.Server
	logger     *zap.Logger
	storage    storage.Storage
}

func New(config *Config) (*APIServer, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %s", err)
	}

	return &APIServer{
		config: config,
		logger: logger,
	}, nil
}

func (apiServer *APIServer) Run() {
	apiServer.logger.Info("starting API server", zap.String("version", version.Version))

	apiServer.startOpenTelemetrySDK()
	apiServer.connectToStorageBackend()
	apiServer.startHTTPServer()

	apiServer.logger.Info("started API server")

	apiServer.waitSignals()
}

func (apiServer *APIServer) startOpenTelemetrySDK() {
	apiServer.logger.Info("starting OpenTelemetry SDK")

	if err := tracing.Start(context.Background(), apiServer.config.Tracing); err != nil {
		apiServer.logger.Error("failed to start OpenTelemetry SDK", zap.Error(err))

		os.Exit(1)
	}

	apiServer.logger.Info("started OpenTelemetry SDK")
}

func (apiServer *APIServer) connectToStorageBackend() {
	if apiServer.storage != nil {
		return
	}

	apiServer.logger.Info("connecting to storage backend", zap.String("backend", apiServer.config.Storage.Backend))

	var err error

	apiServer.storage, err = storagefactory.Create(apiServer.config.Storage)

	if err != nil {
		apiServer.logger.Error("failed to create storage backend client", zap.Error(err))

		os.Exit(1)
	}

	if err = apiServer.storage.Connect(context.Background()); err != nil {
		apiServer.logger.Error("failed to connect storage backend", zap.Error(err))

		os.Exit(1)
	}

	apiServer.logger.Info("connected to storage backend")
}

func (apiServer *APIServer) startHTTPServer() {
	if apiServer.httpServer != nil {
		return
	}

	apiServer.logger.Info("starting HTTP server", zap.String("address", apiServer.config.ListenAddress))

	mux := http.NewServeMux()

	rest.RegisterRoutes(mux, apiServer.storage)

	apiServer.httpServer = &http.Server{
		Addr:              apiServer.config.ListenAddress,
		Handler:           mux,
		ReadTimeout:       apiServer.config.ReadTimeout,
		ReadHeaderTimeout: apiServer.config.ReadHeaderTimeout,
		IdleTimeout:       apiServer.config.IdleTimeout,
		WriteTimeout:      apiServer.config.WriteTimeout,
	}

	go func() {
		if err := apiServer.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			apiServer.logger.Error("failed to start HTTP server", zap.Error(err))

			os.Exit(1)
		}
	}()

	apiServer.logger.Info("started HTTP server")
}

func (apiServer *APIServer) waitSignals() {
	signalChannel := make(chan os.Signal, 2)

	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	receivedSignal := <-signalChannel

	defer apiServer.logger.Sync()

	apiServer.logger.Info("received signal", zap.String("signal", receivedSignal.String()))

	isForceShutdown := receivedSignal == syscall.SIGTERM

	apiServer.stopHTTPServer(isForceShutdown)
	apiServer.disconnectFromStorageBackend()
	apiServer.stopOpenTelemetrySDK()
}

func (apiServer *APIServer) stopHTTPServer(force bool) {
	if apiServer.httpServer == nil {
		return
	}

	apiServer.logger.Info("stopping HTTP server", zap.Bool("force", force))

	var err error

	if force {
		err = apiServer.httpServer.Close()
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), apiServer.config.GracefulShutdownTimeout)

		defer cancel()

		err = apiServer.httpServer.Shutdown(ctx)
	}

	if err != nil {
		apiServer.logger.Error("failed to stop HTTP server", zap.Error(err))
	} else {
		apiServer.logger.Info("stopped HTTP server")
	}
}

func (apiServer *APIServer) disconnectFromStorageBackend() {
	if apiServer.storage == nil {
		return
	}

	apiServer.logger.Info("disconnecting from storage backend")

	if err := apiServer.storage.Disconnect(context.Background()); err != nil {
		apiServer.logger.Error("failed to Disconnect from storage backend")
	} else {
		apiServer.logger.Info("disconnected from storage backend")
	}
}

func (apiServer *APIServer) stopOpenTelemetrySDK() {
	apiServer.logger.Info("stopping OpenTelemetry SDK")

	if err := tracing.Shutdown(context.Background()); err != nil {
		apiServer.logger.Info("failed to stop OpenTelemetry SDK", zap.Error(err))
	} else {
		apiServer.logger.Info("stopped OpenTelemetry SDK")
	}
}
