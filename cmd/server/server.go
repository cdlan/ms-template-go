package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"cdlab.cdlan.net/cdlan/uservices/ms-template/internal/config"
	"cdlab.cdlan.net/cdlan/uservices/ms-template/internal/database"
	grpcimpl "cdlab.cdlan.net/cdlan/uservices/ms-template/internal/grpc"
	"cdlab.cdlan.net/cdlan/uservices/ms-template/internal/grpc/gen"
)

const (
	serviceName string = "misty-products"
	version     string = "[VERSION]" // will be overwritten in pipeline
)

func init() {

	// load config
	config.LoadConfiguration()

	// otel
	if config.C.Otel.Enabled {

		// exporter
		exporter, err := config.C.Otel.NewTracerProvider()
		if err != nil {
			panic(err)
		}

		res := config.C.Otel.NewResource(serviceName, version)
		config.C.Otel.InitTracerProvider(res, exporter)
	}

	// connect to DB
	database.Init()
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
	if config.C.Debug {
		log.Printf("[%s:%s] Started Listener on %s", serviceName, version, listenAddress)
	}

	var opts []grpc.ServerOption

	opts = append(opts, grpc.UnaryInterceptor(config.C.Otel.UnaryServerInterceptor()))
	opts = append(opts, grpc.StreamInterceptor(config.C.Otel.StreamServerInterceptor()))

	grpcServer := grpc.NewServer(opts...)

	//TODO:  register servers
	gen.RegisterCoverageServiceServer(grpcServer, grpcimpl.NewCoverageServerServer())

	reflection.Register(grpcServer)

	// start listening
	if err := grpcServer.Serve(lis); err != nil {

		log.Fatalf("Failed to serve grpc: %v", err)
	}
}
