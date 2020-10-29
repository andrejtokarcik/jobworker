package mtls_test

import (
	"context"
	"net"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/andrejtokarcik/jobworker/client"
	"github.com/andrejtokarcik/jobworker/server"
)

type BufconnConfig struct {
	BufSize       int
	ClientTimeout time.Duration
}

type BufconnSuite struct {
	suite.Suite
	BufconnConfig
	grpcServer *grpc.Server
	listener   *bufconn.Listener
}

func NewBufconnSuite() (suite BufconnSuite) {
	suite.BufconnConfig = BufconnConfig{
		BufSize:       1024 * 1024,
		ClientTimeout: 1 * time.Second,
	}
	return
}

func (suite *BufconnSuite) SetupBufconn(opts ...grpc.ServerOption) {
	suite.grpcServer = server.New(opts...)
	suite.listener = bufconn.Listen(suite.BufSize)
	go func() {
		if err := suite.grpcServer.Serve(suite.listener); err != nil {
			panic(err)
		}
	}()
}

func (suite *BufconnSuite) TearDownBufconn() {
	suite.listener.Close()
	suite.grpcServer.Stop()
}

func (suite *BufconnSuite) contextDialer(context.Context, string) (net.Conn, error) {
	return suite.listener.Dial()
}

func (suite *BufconnSuite) DialBufconn(serverName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return client.DialContextWithTimeout(
		context.Background(),
		suite.ClientTimeout,
		serverName,
		append(
			opts,
			grpc.WithContextDialer(suite.contextDialer),
		)...,
	)
}
