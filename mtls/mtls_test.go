package mtls_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/andrejtokarcik/jobworker/mtls"
)

type mTLSTestSuite struct {
	BufconnSuite
}

type mTLSTestCase struct {
	clientCredsFiles mtls.CredsFiles
	serverName       string
	expectedErr      error
}

func (suite *mTLSTestSuite) SetupSuite() {
	serverCreds, err := mtls.NewServerCreds(mtls.CredsFiles{
		Cert:       "testdata/server-ca/server1.crt",
		Key:        "testdata/server-ca/server1.key",
		PeerCACert: "testdata/client-ca.crt",
	})
	if err != nil {
		panic(err)
	}

	suite.SetupBufconn(grpc.Creds(serverCreds))
}

func (suite *mTLSTestSuite) TearDownSuite() {
	suite.TearDownBufconn()
}

func (suite *mTLSTestSuite) runTestCase(tc mTLSTestCase) {
	clientCreds, err := mtls.NewClientCreds(tc.clientCredsFiles)
	suite.Require().Nil(err, err)

	conn, err := suite.DialBufconn(
		tc.serverName,
		grpc.WithTransportCredentials(clientCreds),
	)
	if conn != nil {
		defer conn.Close()
	}

	if tc.expectedErr == nil {
		suite.Require().Nil(err, err)
	} else {
		suite.Require().NotNil(err)
		suite.Assert().Contains(err.Error(), tc.expectedErr.Error())
	}
}

func validTestCase() mTLSTestCase {
	return mTLSTestCase{
		clientCredsFiles: mtls.CredsFiles{
			Cert:       "testdata/client-ca/client1.crt",
			Key:        "testdata/client-ca/client1.key",
			PeerCACert: "testdata/server-ca.crt",
		},
		serverName:  "server1",
		expectedErr: nil,
	}
}

func (suite *mTLSTestSuite) TestValidCreds() {
	tc := validTestCase()
	suite.runTestCase(tc)
}

func (suite *mTLSTestSuite) TestWrongServerCA() {
	tc := validTestCase()
	tc.clientCredsFiles.PeerCACert = "testdata/server-ca2.crt"
	tc.expectedErr = errors.New("x509: certificate signed by unknown authority")
	suite.runTestCase(tc)
}

func (suite *mTLSTestSuite) TestSelfSignedClientCert() {
	tc := validTestCase()
	tc.clientCredsFiles.Cert = "testdata/self-signed.crt"
	tc.clientCredsFiles.Key = "testdata/self-signed.key"
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
	suite.Run(t, &mTLSTestSuite{NewBufconnSuite()})
}
