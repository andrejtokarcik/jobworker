package server

import (
	"google.golang.org/grpc"
)

func New(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(
		append(
			opts,
			grpc.UnaryInterceptor(AttachClientSubject),
		)...,
	)
}
