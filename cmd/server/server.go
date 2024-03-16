package main

import (
	"context"
	"log"
	"net"

	"ms-template-go/internal/config"
	"ms-template-go/internal/database"
	grpcimpl "ms-template-go/internal/grpc"
	"ms-template-go/internal/grpc/gen"
	"ms-template-go/internal/migrations"
	prom "ms-template-go/pkg/prometheus"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	serviceName string = "misty-products"
	version     string = "[VERSION]" // will be overwritten in pipeline
)

func init() {

	// load config
	config.LoadConfiguration()

	// migrate
	migrations.Migrate(config.C.DB.MigrationPath, config.C.DB.GetURL(), migrations.UP)

	// init DB
	config.C.DB.Init(config.C.Otel.Enabled)

	// init otel
	err := config.C.Otel.Init(serviceName, version)
	if err != nil {

		log.Println(err)
		config.C.Otel.Enabled = false
	}

	// prometheus config
	prom.Conf = prom.Config{
		Port: 2112,
		Path: "/metrics",
	}

	// start exposing metrics
	prom.Conf.ExposeMetrics()
}

func main() {

	// close DB
	defer database.Close()

	// close OTEL
	defer func() {
		if config.C.Otel.Enabled {
			if err := config.C.Otel.ShutdownTracerProvider(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}
	}()

	// open tcp listener
	listenAddress := config.C.GetListenAddress()
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {

		log.Fatalf("failed to listen: %v", err)
	}

	// log
	log.Printf("[%s:%s] Started Listener on %s", serviceName, version, listenAddress)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// register servers
	healthgrpc.RegisterHealthServer(grpcServer, health.NewServer())
	gen.RegisterExampleServer(grpcServer, grpcimpl.NewExampleServer())

	// add reflections
	reflection.Register(grpcServer)

	// start listening
	if err := grpcServer.Serve(lis); err != nil {

		log.Fatalf("Failed to serve grpc: %v", err)
	}
}
