package server

import (
	"context"
	"crypto/x509/pkix"

	pb "github.com/andrejtokarcik/jobworker/proto"
)

func NewCustomJobWorkerServer(cmdCreator CmdCreator) pb.JobWorkerServer {
	return &jobWorkerServer{
		cmdCreator: cmdCreator,
	}
}

func SetClientSubject(ctx context.Context, clientSubject pkix.Name) context.Context {
	return context.WithValue(ctx, clientSubjectKey{}, clientSubject)
}
