package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"net/url"
	"object-t.com/hackz-giganoto/pkg/telemetry"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"goa.design/clue/debug"
	"goa.design/clue/log"
	profileapi "object-t.com/hackz-giganoto/microservices/profile"
	profile "object-t.com/hackz-giganoto/microservices/profile/gen/profile"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		grpcPortF = flag.String("grpc-port", "", "gRPC port (overrides host gRPC port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	// Initialize OpenTelemetry
	config := telemetry.DefaultConfig("profile-service")
	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
		// Remove http:// prefix if present for gRPC connection
		if strings.HasPrefix(endpoint, "http://") {
			endpoint = strings.TrimPrefix(endpoint, "http://")
		}
		config.CollectorAddr = endpoint
	}
	if serviceName := os.Getenv("OTEL_SERVICE_NAME"); serviceName != "" {
		config.ServiceName = serviceName
	}
	if serviceVersion := os.Getenv("OTEL_SERVICE_VERSION"); serviceVersion != "" {
		config.ServiceVersion = serviceVersion
	}

	_, cleanup, err := telemetry.Init(context.Background(), config)
	if err != nil {
		log.Fatalf(context.Background(), err, "failed to initialize OpenTelemetry")
	}
	defer cleanup()

	// Setup logger. Replace logger with your own log package of choice.
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format))
	if *dbgF {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Initialize the services.
	var (
		profileSvc profile.Service
	)
	{
		profileSvc = profileapi.NewProfile(redisClient)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		profileEndpoints *profile.Endpoints
	)
	{
		profileEndpoints = profile.NewEndpoints(profileSvc)
		profileEndpoints.Use(debug.LogPayloads())
		profileEndpoints.Use(log.Endpoint)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "localhost":
		{
			addr := "grpc://localhost:50052"
			u, err := url.Parse(addr)
			if err != nil {
				log.Fatalf(ctx, err, "invalid URL %#v\n", addr)
			}
			if *secureF {
				u.Scheme = "grpcs"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *grpcPortF != "" {
				h, _, err := net.SplitHostPort(u.Host)
				if err != nil {
					log.Fatalf(ctx, err, "invalid URL %#v\n", u.Host)
				}
				u.Host = net.JoinHostPort(h, *grpcPortF)
			} else if u.Port() == "" {
				u.Host = net.JoinHostPort(u.Host, "50052")
			}
			handleGRPCServer(ctx, u, profileEndpoints, &wg, errc, *dbgF)
		}

	default:
		log.Fatal(ctx, fmt.Errorf("invalid host argument: %q (valid hosts: localhost)", *hostF))
	}

	// Wait for signal.
	log.Printf(ctx, "exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	log.Printf(ctx, "exited")
}
