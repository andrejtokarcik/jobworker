package server_test

import (
	"context"
	"crypto/x509/pkix"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	pb "github.com/andrejtokarcik/jobworker/proto"
	"github.com/andrejtokarcik/jobworker/server"
	"github.com/andrejtokarcik/jobworker/server/mocks"
)

type jobWorkerTestSuite struct {
	suite.Suite
	mockCmdCreator *mocks.CmdCreator
	jobWorker      pb.JobWorkerServer
	ctx            context.Context
}

func (suite *jobWorkerTestSuite) SetupTest() {
	suite.mockCmdCreator = new(mocks.CmdCreator)
	suite.jobWorker = server.NewCustomJobWorkerServer(suite.mockCmdCreator)

	clientSubject := pkix.Name{CommonName: "client"}
	suite.ctx = server.SetClientSubject(context.Background(), clientSubject)
}

func (suite *jobWorkerTestSuite) TearDownTest() {
	suite.mockCmdCreator.AssertExpectations(suite.T())
}

func (suite *jobWorkerTestSuite) TestFailedState() {
	mockCmd, startResp, err := suite.callStartJob()
	suite.Require().Nil(err, err)

	getResp, err := suite.callGetJob(mockCmd, startResp.JobUUID, cmdStatusTerminated())
	suite.Require().Nil(err, err)
	suite.Assert().Equal(getResp.State, pb.GetJobResponse_FAILED, getResp.State)
}

func (suite *jobWorkerTestSuite) TestStoppedState() {
	mockCmd, startResp, err := suite.callStartJob()
	suite.Require().Nil(err, err)

	err = suite.callStopJob(mockCmd, startResp.JobUUID, cmdStatusRunning())
	suite.Require().Nil(err, err)

	getResp, err := suite.callGetJob(mockCmd, startResp.JobUUID, cmdStatusTerminated())
	suite.Require().Nil(err, err)
	suite.Assert().Equal(getResp.State, pb.GetJobResponse_STOPPED, getResp.State)
}

func (suite *jobWorkerTestSuite) TestRemainsFailedAfterStop() {
	mockCmd, startResp, err := suite.callStartJob()
	suite.Require().Nil(err, err)

	getResp, err := suite.callGetJob(mockCmd, startResp.JobUUID, cmdStatusTerminated())
	suite.Require().Nil(err, err)
	suite.Assert().Equal(getResp.State, pb.GetJobResponse_FAILED, getResp.State)

	err = suite.callStopJob(mockCmd, startResp.JobUUID, cmdStatusTerminated())
	suite.Assert().NotNil(err)
	suite.Assert().Contains(err.Error(), "not running")

	getResp, err = suite.callGetJob(mockCmd, startResp.JobUUID, cmdStatusTerminated())
	suite.Require().Nil(err, err)
	suite.Assert().Equal(getResp.State, pb.GetJobResponse_FAILED, getResp.State)
}

func (suite *jobWorkerTestSuite) callStartJob() (*mocks.Cmd, *pb.StartJobResponse, error) {
	mockCmd := new(mocks.Cmd)
	suite.mockCmdCreator.On("NewCmd", cmdSpec()).Return(mockCmd).Once()
	mockCmd.On("Start").Return(make(<-chan server.CmdStatus)).Once()
	defer suite.checkAndClearMockCalls(mockCmd)

	startReq := &pb.StartJobRequest{Command: cmdSpec()}
	startResp, err := suite.jobWorker.StartJob(suite.ctx, startReq)
	return mockCmd, startResp, err
}

func (suite *jobWorkerTestSuite) callStopJob(mockCmd *mocks.Cmd, jobUUID string, status server.CmdStatus) error {
	mockCmd.On("Status").Return(status).Once()
	mockCmd.On("Stop").Return(nil).Maybe()
	stopReq := &pb.StopJobRequest{JobUUID: jobUUID}
	defer suite.checkAndClearMockCalls(mockCmd)

	_, err := suite.jobWorker.StopJob(suite.ctx, stopReq)
	return err
}

func (suite *jobWorkerTestSuite) callGetJob(mockCmd *mocks.Cmd, jobUUID string, status server.CmdStatus) (*pb.GetJobResponse, error) {
	mockCmd.On("Spec").Return(cmdSpec()).Once()
	mockCmd.On("Status").Return(status).Once()
	defer suite.checkAndClearMockCalls(mockCmd)

	getReq := &pb.GetJobRequest{JobUUID: jobUUID, WithLogs: false}
	return suite.jobWorker.GetJob(suite.ctx, getReq)
}

func (suite *jobWorkerTestSuite) checkAndClearMockCalls(mockCmd *mocks.Cmd) {
	mockCmd.AssertExpectations(suite.T())
	mockCmd.ExpectedCalls = nil
}

func cmdSpec() *server.CmdSpec {
	return &server.CmdSpec{
		Command: "some-command",
		Args:    []string{"--some-arg", "some-value"},
		Dir:     "/some-path",
		Env:     []string{"X=y"},
	}
}

func cmdStatusRunning() server.CmdStatus {
	return server.CmdStatus{
		PID:      3000,
		Complete: false,
		StartTs:  123,
	}
}

func cmdStatusTerminated() server.CmdStatus {
	running := cmdStatusRunning()
	running.StopTs = running.StartTs + 456
	running.Error = errors.New("signal: terminated")
	return running
}

func TestJobWorkerServer(t *testing.T) {
	suite.Run(t, &jobWorkerTestSuite{})
}
