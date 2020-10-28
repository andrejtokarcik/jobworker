package server

import (
	"context"
	"crypto/x509/pkix"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/andrejtokarcik/jobworker/proto"
)

type JobWorkerServer struct {
	pb.UnimplementedJobWorkerServer
	cmdCreator CmdCreator
	jobs       sync.Map
}

func NewJobWorkerServer() *JobWorkerServer {
	return &JobWorkerServer{
		cmdCreator: gocmdCreator{},
	}
}

type job struct {
	Cmd
	subjectName string
}

func (server *JobWorkerServer) StartJob(ctx context.Context, req *pb.StartJobRequest) (*pb.StartJobResponse, error) {
	jobUuid, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate a new UUID: %v", err)
	}

	cmd := server.cmdCreator.NewCmd(
		req.Command.Dir, req.Command.Env, req.Command.Command, req.Command.Args,
	)
	subjectName := ctx.Value(clientSubjectKey{}).(pkix.Name).CommonName

	server.jobs.Store(jobUuid, job{cmd, subjectName})
	cmd.Start()

	response := &pb.StartJobResponse{
		JobUuid: jobUuid.String(),
	}
	return response, nil
}
