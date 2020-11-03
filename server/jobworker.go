package server

import (
	"context"
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
	startedBy ClientSubject
	stopped   bool
}

func (j job) isOwnedBy(client ClientSubject) bool {
	return j.startedBy.CommonName == client.CommonName
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

	clientSubject, err := getClientSubject(ctx)
	if err != nil {
		return nil, err
	}

	job := job{
		Cmd:       server.cmdCreator.NewCmd(req.Command),
		startedBy: clientSubject,
		stopped:   false,
	}
	server.jobs.Store(jobUUID, job)
	job.Start()

	response := &pb.StartJobResponse{
		JobUUID: jobUUID,
	}
	return response, nil
}

func (server *jobWorkerServer) StopJob(ctx context.Context, req *pb.StopJobRequest) (*pb.StopJobResponse, error) {
	clientSubject, err := getClientSubject(ctx)
	if err != nil {
		return nil, err
	}

	value, ok := server.jobs.Load(req.JobUUID)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no job found for the given UUID")
	}
	job := value.(job)

	if !job.isOwnedBy(clientSubject) {
		return nil, status.Errorf(codes.PermissionDenied, "not allowed to stop this job")
	}

	if state := determineState(job.stopped, job.Status()); state != pb.GetJobResponse_RUNNING {
		return nil, status.Errorf(codes.FailedPrecondition, "job process is not running: %v", state)
	}

	err = job.Stop()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send a stop signal to the job: %v", err)
	}

	job.stopped = true
	server.jobs.Store(req.JobUUID, job)
	return &pb.StopJobResponse{}, nil
}

func (server *jobWorkerServer) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	value, ok := server.jobs.Load(req.JobUUID)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no job found for the given UUID")
	}

	job := value.(job)
	cmdStatus := job.Status()

	startedAt, err := unixNanoToTimestamp(cmdStatus.StartTs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert job start to timestamp: %v", err)
	}

	endedAt, err := unixNanoToTimestamp(cmdStatus.StopTs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert job end to timestamp: %v", err)
	}

	resp := &pb.GetJobResponse{
		Command:   job.Spec(),
		State:     determineState(job.stopped, cmdStatus),
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
