package server

import (
	"google.golang.org/grpc"

	pb "github.com/andrejtokarcik/jobworker/proto"
)

func New(filter RPCCallerFilter, opts ...grpc.ServerOption) *grpc.Server {
	server := grpc.NewServer(
		append(
			opts,
			grpc.ChainUnaryInterceptor(
				AttachClientSubject,
				ApplyRPCCallerFilter(filter),
			),
		)...,
	)
	pb.RegisterJobWorkerServer(server, NewJobWorkerServer())
	return server
}
