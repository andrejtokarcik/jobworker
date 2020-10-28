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
	Start() <-chan CmdStatus
	Status() CmdStatus
	Stop() error
	Spec() *CmdSpec
	State() CmdState
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

func (cmd goCmd) Spec() (spec *CmdSpec) {
	spec.Command = cmd.Name
	spec.Args = cmd.Args
	spec.Env = cmd.Env
	spec.Dir = cmd.Dir
	return
}

func (cmd goCmd) State() CmdState {
	status := cmd.Status()

	if status.Error != nil {
		return pb.GetJobResponse_FAILED
	}

	if status.StartTs == 0 {
		return pb.GetJobResponse_PENDING
	}

	if status.StopTs == 0 {
		return pb.GetJobResponse_RUNNING
	}

	if status.Complete {
		return pb.GetJobResponse_COMPLETED
	} else {
		return pb.GetJobResponse_STOPPED
	}
}
