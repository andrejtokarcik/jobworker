package mtls_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/andrejtokarcik/jobworker/mtls"
	"github.com/andrejtokarcik/jobworker/test"
	"github.com/andrejtokarcik/jobworker/test/data"
)

type mTLSTestSuite struct {
	test.BufconnSuite
}

type mTLSTestCase struct {
	clientCredsFiles mtls.CredsFiles
	serverName string
	expectedErr      error
}

func (suite *mTLSTestSuite) SetupSuite() {
	grpcServer := grpc.NewServer(
		grpc.Creds(testdata.DefaultServerCreds()),
	)
	suite.SetupBufconn(grpcServer)
}

func (suite *mTLSTestSuite) TearDownSuite() {
	suite.TearDownBufconn()
}

func (suite *mTLSTestSuite) runTestCase(tc mTLSTestCase) {
	clientCreds, err := mtls.NewClientCreds(
		testdata.CredsFilesPaths(tc.clientCredsFiles),
	)
	suite.Require().Nil(err)

	conn, err := suite.DialBufconn(tc.serverName, clientCreds)
	if conn != nil {
		defer conn.Close()
	}

	if tc.expectedErr == nil {
		suite.Require().Nil(err)
	} else {
		suite.Require().NotNil(err)
		suite.Assert().Contains(err.Error(), tc.expectedErr.Error())
	}
}

func validTestCase() mTLSTestCase {
	return mTLSTestCase{
		clientCredsFiles: testdata.DefaultClientCredsFiles(),
		serverName: "server1",
		expectedErr:      nil,
	}
}

func (suite *mTLSTestSuite) TestValidCreds() {
	tc := validTestCase()
	suite.runTestCase(tc)
}

func (suite *mTLSTestSuite) TestWrongServerCA() {
	tc := validTestCase()
	tc.clientCredsFiles.PeerCACert = "server-ca2.crt"
	tc.expectedErr = errors.New("x509: certificate signed by unknown authority")
	suite.runTestCase(tc)
}

func (suite *mTLSTestSuite) TestSelfSignedClientCert() {
	tc := validTestCase()
	tc.clientCredsFiles.Cert = "self-signed.crt"
	tc.clientCredsFiles.Key = "self-signed.key"
	tc.expectedErr = errors.New("context deadline exceeded")
	suite.runTestCase(tc)
}

func (suite *mTLSTestSuite) TestInvalidServerName() {
	tc := validTestCase()
	tc.serverName = "server2"
	tc.expectedErr = errors.New("x509: certificate is valid for server1, not server2")
	suite.runTestCase(tc)
}

func TestMutualTLS(t *testing.T) {
	suite.Run(t, &mTLSTestSuite{test.NewBufconnSuite()})
}
