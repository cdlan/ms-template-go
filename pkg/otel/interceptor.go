package otel

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
)

func (c *Config) UnaryServerInterceptor() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		if !c.Enabled {
			return handler(ctx, req)
		}

		traceName := info.FullMethod
		newCtx, span := otel.Tracer(traceName).Start(ctx, info.FullMethod)
		defer span.End()

		// save trace name in context
		newCtx = context.WithValue(newCtx, "trace_name", traceName)

		resp, err = handler(newCtx, req)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return resp, err
	}
}
func (c *Config) StreamServerInterceptor() grpc.StreamServerInterceptor {

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		if !c.Enabled {
			return handler(srv, stream)
		}

		traceName := info.FullMethod

		newCtx, span := otel.Tracer(traceName).Start(stream.Context(), info.FullMethod)
		defer span.End()

		// save trace name in context
		newCtx = context.WithValue(newCtx, "trace_name", traceName)

		err := handler(srv, &wrappedStream{stream, newCtx})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return err
	}
}
func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) RecvMsg(m interface{}) error {

	name := fmt.Sprintf("grpc-received: %T - %v\n", m, m)

	newCtx, span := otel.Tracer(name).Start(w.Context(), "RecvMsg")
	defer span.End()
	w.ctx = newCtx

	return w.ServerStream.RecvMsg(m)
}
func (w *wrappedStream) SendMsg(m interface{}) error {

	name := fmt.Sprintf("grpc-received: %T - %v\n", m, m)

	newCtx, span := otel.Tracer(name).Start(w.Context(), "SendMsg")
	defer span.End()
	w.ctx = newCtx

	return w.ServerStream.SendMsg(m)
}
