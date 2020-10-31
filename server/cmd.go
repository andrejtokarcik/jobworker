package server

import (
	gocmd "github.com/go-cmd/cmd"

	pb "github.com/andrejtokarcik/jobworker/proto"
)

type CmdStatus = gocmd.Status
type CmdSpec = pb.CommandSpec

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
