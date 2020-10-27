package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func Dial(target string, timeout time.Duration, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return grpc.DialContext(
		ctx,
		target,
		append(
			opts,
			grpc.WithReturnConnectionError(),
		)...,
	)
}
