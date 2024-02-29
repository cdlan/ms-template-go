package grpc

import (
	"context"
	"fmt"

	"ms-template-go/internal/config"
	"ms-template-go/internal/grpc/gen"
)

type ExampleServer struct {
	gen.UnimplementedExampleServer
}

func NewExampleServer() *ExampleServer {
	return &ExampleServer{}
}

func (S *ExampleServer) StreamMethod(req *gen.GetAll, stream gen.Example_StreamMethodServer) error {

	// start span
	ctx, span := config.C.Otel.NewSpan(stream.Context(), "main.Example/StreamMethod")
	if span != nil {
		defer span.End()
	}

	fmt.Println(ctx) // it does nothing but throws error otherwise
	
	// do stuff

	return nil
}
func (S *ExampleServer) UnaryMethod(ctx context.Context, req *gen.Get) (*gen.Item, error) {

	// start span
	ctx, span := config.C.Otel.NewSpan(ctx, "main.Example/UnaryMethod")
	if span != nil {
		defer span.End()
	}

	// do stuff

	return &gen.Item{}, nil
}