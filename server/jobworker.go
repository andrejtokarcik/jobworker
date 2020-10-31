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
	stopped    bool
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

	clientSubject := ctx.Value(clientSubjectKey{})
	if clientSubject == nil {
		return nil, status.Errorf(codes.Internal, "cannot determine client subject")
	}

	job := job{
		Cmd:        server.cmdCreator.NewCmd(req.Command),
		clientName: clientSubject.(pkix.Name).CommonName,
		stopped:    false,
	}
	server.jobs.Store(jobUUID, job)
	job.Start()

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

	job := value.(job)
	if _, state := job.Status(); state != pb.GetJobResponse_RUNNING {
		return nil, status.Errorf(codes.FailedPrecondition, "job process is not running: %v", state)
	}

	err := job.Stop()
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
	jobStatus, jobState := job.Status()

	startedAt, err := unixNanoToTimestamp(jobStatus.StartTs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert job start to timestamp: %v", err)
	}

	endedAt, err := unixNanoToTimestamp(jobStatus.StopTs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert job end to timestamp: %v", err)
	}

	resp := &pb.GetJobResponse{
		Command:   job.Spec(),
		State:     jobState,
		ExitCode:  int32(jobStatus.Exit),
		StartedAt: startedAt,
		EndedAt:   endedAt,
		PID:       uint32(jobStatus.PID),
	}

	if req.WithLogs {
		resp.Stdout = jobStatus.Stdout
		resp.Stderr = jobStatus.Stderr
	}

	return resp, nil
}

func (j job) Status() (CmdStatus, pb.GetJobResponse_State) {
	status := j.Cmd.Status()

	if j.stopped {
		return status, pb.GetJobResponse_STOPPED
	}

	if status.Error != nil {
		return status, pb.GetJobResponse_FAILED
	}

	if status.StartTs == 0 {
		return status, pb.GetJobResponse_PENDING
	}

	if status.StopTs == 0 {
		return status, pb.GetJobResponse_RUNNING
	}

	if status.Complete {
		return status, pb.GetJobResponse_COMPLETED
	}

	return status, pb.GetJobResponse_UNKNOWN
}

func unixNanoToTimestamp(n int64) (*types.Timestamp, error) {
	return types.TimestampProto(time.Unix(0, n))
}
