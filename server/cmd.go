package server

import (
	gocmd "github.com/go-cmd/cmd"
)

type CmdCreator interface {
	NewCmd(dir string, env []string, name string, args []string) Cmd
}

type Cmd interface {
	Start() <-chan CmdStatus
	Status() CmdStatus
	Stop() error
}

type CmdStatus = gocmd.Status

type gocmdCreator struct{}

func (gocmdCreator) NewCmd(dir string, env []string, name string, args []string) Cmd {
	c := gocmd.NewCmd(name, args...)
	c.Dir = dir
	c.Env = env
	return c
}
