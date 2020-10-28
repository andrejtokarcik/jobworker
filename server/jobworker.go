package server

import (
	"context"
	"crypto/x509/pkix"
	"sync"
	"time"

	"github.com/gogo/protobuf/types"
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
		cmdCreator: goCmdCreator{},
	}
}

type job struct {
	Cmd
	clientName string
}

func (server *JobWorkerServer) StartJob(ctx context.Context, req *pb.StartJobRequest) (*pb.StartJobResponse, error) {
	jobUuid, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate a new UUID: %v", err)
	}

	cmd := server.cmdCreator.NewCmd(req.Command)

	clientSubject := ctx.Value(clientSubjectKey{})
	if clientSubject == nil {
		return nil, status.Errorf(codes.Internal, "cannot determine client subject")
	}
	clientName := clientSubject.(pkix.Name).CommonName

	server.jobs.Store(jobUuid, job{cmd, clientName})
	cmd.Start()

	response := &pb.StartJobResponse{
		JobUuid: jobUuid.String(),
	}
	return response, nil
}

func (server *JobWorkerServer) StopJob(ctx context.Context, req *pb.StopJobRequest) (*pb.StopJobResponse, error) {
	value, ok := server.jobs.Load(req.JobUuid)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no job found for the given UUID")
	}

	// Stop's error is irrelevant: even when it occurs, it need not be handled
	_ = value.(job).Stop()

	return &pb.StopJobResponse{}, nil
}

func (server *JobWorkerServer) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	value, ok := server.jobs.Load(req.JobUuid)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no job found for the given UUID")
	}

	cmd := value.(job).Cmd

	startedAt, err := unixNanoToTimestamp(cmd.Status().StartTs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert job start to timestamp: %v", err)
	}

	endedAt, err := unixNanoToTimestamp(cmd.Status().StopTs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert job end to timestamp: %v", err)
	}

	resp := &pb.GetJobResponse{
		Command:   cmd.Spec(),
		State:     cmd.State(),
		ExitCode:  int32(cmd.Status().Exit),
		StartedAt: startedAt,
		EndedAt:   endedAt,
		PID:       uint32(cmd.Status().PID),
	}

	if req.WithLogs {
		resp.Stdout = cmd.Status().Stdout
		resp.Stderr = cmd.Status().Stderr
	}

	return resp, nil
}

func unixNanoToTimestamp(n int64) (*types.Timestamp, error) {
	return types.TimestampProto(time.Unix(0, n))
}
