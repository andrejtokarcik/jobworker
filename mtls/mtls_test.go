package mtls_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/andrejtokarcik/jobworker/mtls"
	"github.com/andrejtokarcik/jobworker/test"
	"github.com/andrejtokarcik/jobworker/test/data"
)

type mTLSTestSuite struct {
	test.BufconnSuite
}

type mTLSTestCase struct {
	clientCredsFiles mtls.CredsFiles
	expectedErr      error
}

func (suite *mTLSTestSuite) SetupSuite() {
	suite.SetupBufconnWithDefaultCreds()
}

func (suite *mTLSTestSuite) TearDownSuite() {
	suite.TearDownBufconn()
}

func (suite *mTLSTestSuite) runTestCase(tc mTLSTestCase) {
	clientCreds, err := mtls.NewClientCreds(
		testdata.CredsFilesPaths(tc.clientCredsFiles),
	)
	suite.Require().Nil(err)

	conn, err := suite.DialBufconn(clientCreds)
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

func TestMutualTLS(t *testing.T) {
	suite.Run(t, &mTLSTestSuite{test.NewBufconnSuite()})
}
