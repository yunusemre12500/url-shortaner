package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.yunus-emre.dev/url-shortaner/pkg/version"
)

func Start(ctx context.Context, config *Config) error {
	res, err := newResource(ctx, config)

	if err != nil {
		return err
	}

	spanExporter, err := newSpanExporter(ctx, config)

	if err != nil {
		return err
	}

	tp := newTracerProvider(config, res, spanExporter)

	otel.SetTracerProvider(tp)

	return nil
}

func Shutdown(ctx context.Context) error {
	if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		return tp.Shutdown(ctx)
	}

	return nil
}

func newSpanExporter(ctx context.Context, config *Config) (trace.SpanExporter, error) {
	var (
		err          error
		spanExporter trace.SpanExporter
	)

	switch config.Protocol {
	case "grpc":
		opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(config.Endpoint)}

		if config.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}

		spanExporter, err = otlptracegrpc.New(ctx, opts...)

		if err != nil {
			return nil, err
		}
	case "http":
		opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(config.Endpoint)}

		if config.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		spanExporter, err = otlptracehttp.New(ctx, opts...)

		if err != nil {
			return nil, err
		}
	}

	return spanExporter, nil
}

func newResource(ctx context.Context, config *Config) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(config.ServiceName), semconv.ServiceVersion(version.Version)),
		resource.WithContainer(),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithService(),
		resource.WithTelemetrySDK())
}

func newTracerProvider(config *Config, res *resource.Resource, spanExporter trace.SpanExporter) *trace.TracerProvider {
	opts := []sdktrace.TracerProviderOption{sdktrace.WithResource(res)}

	if config.Batching.Enabled {
		opts = append(opts, sdktrace.WithBatcher(spanExporter))
	} else {
		opts = append(opts, sdktrace.WithSyncer(spanExporter))
	}

	tp := sdktrace.NewTracerProvider(opts...)

	return tp
}
