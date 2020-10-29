package server

import (
	"google.golang.org/grpc"

	pb "github.com/andrejtokarcik/jobworker/proto"
)

func New(opts ...grpc.ServerOption) *grpc.Server {
	server := grpc.NewServer(
		append(
			opts,
			grpc.UnaryInterceptor(AttachClientSubject),
		)...,
	)
	pb.RegisterJobWorkerServer(server, NewJobWorkerServer())
	return server
}
