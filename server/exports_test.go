package server

import (
	"context"
	"crypto/x509/pkix"
)

func NewCustomJobWorkerServer(cmdCreator CmdCreator) *JobWorkerServer {
	return &JobWorkerServer{
		cmdCreator: cmdCreator,
	}
}

func SetClientSubject(ctx context.Context, clientSubject pkix.Name) context.Context {
	return context.WithValue(ctx, clientSubjectKey{}, clientSubject)
}
