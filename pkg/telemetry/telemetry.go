package telemetry

import (
	"context"
	"crypto/tls"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"goa.design/clue/clue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
)

// Config holds telemetry configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	Region         string
	CollectorAddr  string
	UseInsecure    bool
}

// Init initializes OpenTelemetry for the given service
func Init(ctx context.Context, config Config) (*clue.Config, func(), error) {
	// Create credentials based on configuration
	var creds credentials.TransportCredentials
	if config.UseInsecure {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(&tls.Config{})
	}

	// Create span exporter
	spanExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(config.CollectorAddr),
		otlptracegrpc.WithTLSCredentials(creds),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create span exporter: %w", err)
	}

	// Create metric exporter
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(config.CollectorAddr),
		otlpmetricgrpc.WithTLSCredentials(creds),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	// Note: OpenTelemetry Go logs are experimental and not production-ready
	// Using structured JSON logging to stdout instead

	// Create clue configuration
	cfg, err := clue.NewConfig(ctx,
		config.ServiceName,
		config.ServiceVersion,
		metricExporter,
		spanExporter,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create clue config: %w", err)
	}

	// Configure trace propagation for W3C trace context
	// This ensures proper trace context propagation across services
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	// Configure OpenTelemetry with clue
	clue.ConfigureOpenTelemetry(ctx, cfg)

	// Create cleanup function
	cleanup := func() {
		// Shutdown exporters
		if spanExporter != nil {
			spanExporter.Shutdown(ctx)
		}
		if metricExporter != nil {
			metricExporter.Shutdown(ctx)
		}
	}

	log.Printf("OpenTelemetry initialized for service: %s", config.ServiceName)

	return cfg, cleanup, nil
}

// HTTPMiddleware returns HTTP middleware for OpenTelemetry tracing
func HTTPMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return otelhttp.NewHandler(h, serviceName)
	}
}

// GRPCServerInterceptor returns gRPC server interceptor for OpenTelemetry tracing
func GRPCServerInterceptor() grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler())
}

// GRPCClientInterceptor returns gRPC client interceptor for OpenTelemetry tracing
func GRPCClientInterceptor() grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler())
}

// DefaultConfig returns default telemetry configuration
func DefaultConfig(serviceName string) Config {
	return Config{
		ServiceName:    serviceName,
		ServiceVersion: "1.0.0",
		Environment:    "development",
		Region:         "local",
		CollectorAddr:  "otel-collector:4317",
		UseInsecure:    true,
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp  string            `json:"timestamp"`
	Level      string            `json:"level"`
	Service    string            `json:"service"`
	Message    string            `json:"message"`
	TraceID    string            `json:"trace_id,omitempty"`
	SpanID     string            `json:"span_id,omitempty"`
	Error      string            `json:"error,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// getServiceName extracts service name from environment or default
func getServiceName() string {
	if serviceName := os.Getenv("OTEL_SERVICE_NAME"); serviceName != "" {
		return serviceName
	}
	if serviceName := os.Getenv("SERVICE_NAME"); serviceName != "" {
		return serviceName
	}
	return "unknown-service"
}
