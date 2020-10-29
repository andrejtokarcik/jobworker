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

type serverTestSuite struct {
	suite.Suite
	mockCmdCreator *mocks.CmdCreator
	jobWorker      pb.JobWorkerServer
	ctx            context.Context
}

func (suite *serverTestSuite) SetupTest() {
	suite.mockCmdCreator = new(mocks.CmdCreator)
	suite.jobWorker = server.NewCustomJobWorkerServer(suite.mockCmdCreator)

	clientSubject := pkix.Name{CommonName: "client"}
	suite.ctx = server.SetClientSubject(context.Background(), clientSubject)
}

func (suite *serverTestSuite) TestStartJob() {
	someCmdSpec := &server.CmdSpec{
		Command: "some-command",
		Args:    []string{"--some-arg", "some-value"},
		Dir:     "/some-path",
		Env:     []string{"X=y"},
	}

	req := &pb.StartJobRequest{
		Command: someCmdSpec,
	}

	mockCmd := new(mocks.Cmd)
	suite.mockCmdCreator.On("NewCmd", someCmdSpec).Return(mockCmd).Once()
	mockCmd.On("Start").Return(make(<-chan server.CmdStatus)).Once()

	_, err := suite.jobWorker.StartJob(suite.ctx, req)
	suite.Require().Nil(err, err)

	mockCmd.AssertExpectations(suite.T())
	suite.mockCmdCreator.AssertExpectations(suite.T())
}

func TestJobWorkerServer(t *testing.T) {
	suite.Run(t, &serverTestSuite{})
}
