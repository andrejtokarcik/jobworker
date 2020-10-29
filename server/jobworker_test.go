package server_test

import (
	"context"
	"crypto/x509/pkix"
	"testing"

	"github.com/stretchr/testify/suite"

	pb "github.com/andrejtokarcik/jobworker/proto"
	"github.com/andrejtokarcik/jobworker/server"
	"github.com/andrejtokarcik/jobworker/server/mocks"
)

var (
	someCommand = "some-command"
	someArgs    = []string{"--some-arg", "/some-path", "some-value"}
	someDir     = "."
	someEnv     = []string{"X=y"}
)

type serverTestSuite struct {
	suite.Suite
	mockCmdCreator *mocks.CmdCreator
	jobWorker      *server.JobWorkerServer
	ctx            context.Context
}

func (suite *serverTestSuite) SetupTest() {
	suite.mockCmdCreator = new(mocks.CmdCreator)
	suite.jobWorker = server.NewCustomJobWorkerServer(suite.mockCmdCreator)

	clientSubject := pkix.Name{CommonName: "client"}
	suite.ctx = server.SetClientSubject(context.Background(), clientSubject)
}

func (suite *serverTestSuite) TestStartJob() {
	req := &pb.StartJobRequest{
		Command: &pb.CommandSpec{
			Command: someCommand,
			Args:    someArgs,
			Dir:     someDir,
			Env:     someEnv,
		},
	}

	mockCmd := new(mocks.Cmd)
	suite.mockCmdCreator.On(
		"NewCmd", someDir, someEnv, someCommand, someArgs,
	).Return(mockCmd).Once()
	mockCmd.On("Start").Return(make(<-chan server.CmdStatus)).Once()

	_, err := suite.jobWorker.StartJob(suite.ctx, req)
	suite.Require().Nil(err, err)

	mockCmd.AssertExpectations(suite.T())
	suite.mockCmdCreator.AssertExpectations(suite.T())
}

func TestJobWorkerServer(t *testing.T) {
	suite.Run(t, &serverTestSuite{})
}
