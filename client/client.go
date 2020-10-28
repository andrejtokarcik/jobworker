package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func DialContextWithTimeout(ctx context.Context, timeout time.Duration, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return grpc.DialContext(
		ctxWithTimeout,
		target,
		append(
			opts,
			grpc.WithReturnConnectionError(),
		)...,
	)
}
