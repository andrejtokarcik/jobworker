package server

import (
	gocmd "github.com/go-cmd/cmd"

	pb "github.com/andrejtokarcik/jobworker/proto"
)

type CmdStatus = gocmd.Status
type CmdSpec = pb.CommandSpec
type CmdState = pb.GetJobResponse_State

type CmdCreator interface {
	NewCmd(*CmdSpec) Cmd
}

type Cmd interface {
	Spec() *CmdSpec
	Start() <-chan CmdStatus
	Status() CmdStatus
	Stop() error
}

type goCmdCreator struct{}

type goCmd struct {
	*gocmd.Cmd
}

func (goCmdCreator) NewCmd(spec *CmdSpec) Cmd {
	c := gocmd.NewCmd(spec.Command, spec.Args...)
	c.Dir = spec.Dir
	c.Env = spec.Env
	return goCmd{c}
}

func (cmd goCmd) Spec() *CmdSpec {
	return &CmdSpec{
		Command: cmd.Name,
		Args:    cmd.Args,
		Env:     cmd.Env,
		Dir:     cmd.Dir,
	}
}

func determineState(stopped bool, status CmdStatus) (CmdState, string) {
	if stopped {
		return pb.GetJobResponse_STOPPED, ""
	}

	if status.Error != nil {
		return pb.GetJobResponse_FAILED, status.Error.Error()
	}

	if status.StartTs == 0 {
		return pb.GetJobResponse_PENDING, ""
	}

	if status.StopTs == 0 {
		return pb.GetJobResponse_RUNNING, ""
	}

	if status.Complete {
		return pb.GetJobResponse_COMPLETED, ""
	}

	return pb.GetJobResponse_UNKNOWN, ""
}
