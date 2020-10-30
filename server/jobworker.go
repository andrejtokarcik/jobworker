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

type jobWorkerServer struct {
	pb.UnimplementedJobWorkerServer
	cmdCreator CmdCreator
	jobs       sync.Map
}

func NewJobWorkerServer() pb.JobWorkerServer {
	return &jobWorkerServer{
		cmdCreator: goCmdCreator{},
	}
}

type job struct {
	Cmd
	clientName string
}

func (server *jobWorkerServer) StartJob(ctx context.Context, req *pb.StartJobRequest) (*pb.StartJobResponse, error) {
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate a new UUID: %v", err)
	}
	jobUUID := uuidObj.String()
	if jobUUID == "" {
		return nil, status.Errorf(codes.Internal, "generated UUID is invalid")
	}

	cmd := server.cmdCreator.NewCmd(req.Command)

	clientSubject := ctx.Value(clientSubjectKey{})
	if clientSubject == nil {
		return nil, status.Errorf(codes.Internal, "cannot determine client subject")
	}
	clientName := clientSubject.(pkix.Name).CommonName

	server.jobs.Store(jobUUID, job{cmd, clientName})
	cmd.Start()

	response := &pb.StartJobResponse{
		JobUUID: jobUUID,
	}
	return response, nil
}

func (server *jobWorkerServer) StopJob(ctx context.Context, req *pb.StopJobRequest) (*pb.StopJobResponse, error) {
	value, ok := server.jobs.Load(req.JobUUID)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no job found for the given UUID")
	}

	err := value.(job).Stop()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send a stop signal to the job: %v", err)
	}

	return &pb.StopJobResponse{}, nil
}

func (server *jobWorkerServer) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	value, ok := server.jobs.Load(req.JobUUID)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no job found for the given UUID")
	}

	cmd := value.(job).Cmd
	cmdStatus := cmd.Status()

	startedAt, err := unixNanoToTimestamp(cmdStatus.StartTs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert job start to timestamp: %v", err)
	}

	endedAt, err := unixNanoToTimestamp(cmdStatus.StopTs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert job end to timestamp: %v", err)
	}

	resp := &pb.GetJobResponse{
		Command:   cmd.Spec(),
		State:     determineState(cmdStatus),
		ExitCode:  int32(cmdStatus.Exit),
		StartedAt: startedAt,
		EndedAt:   endedAt,
		PID:       uint32(cmdStatus.PID),
	}

	if req.WithLogs {
		resp.Stdout = cmdStatus.Stdout
		resp.Stderr = cmdStatus.Stderr
	}

	return resp, nil
}

func unixNanoToTimestamp(n int64) (*types.Timestamp, error) {
	return types.TimestampProto(time.Unix(0, n))
}
