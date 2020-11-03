package server

import (
	"context"

	pb "github.com/andrejtokarcik/jobworker/proto"
)

func NewCustomJobWorkerServer(cmdCreator CmdCreator) pb.JobWorkerServer {
	return &jobWorkerServer{
		cmdCreator: cmdCreator,
	}
}

func SetClientSubject(ctx context.Context, clientSubject ClientSubject) context.Context {
	return context.WithValue(ctx, clientSubjectKey{}, clientSubject)
}
