package mtls_test

import (
	"context"
	"net"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
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

func NewBufconnSuite() BufconnSuite {
	return BufconnSuite{
		BufconnConfig: BufconnConfig{
			BufSize:       1024 * 1024,
			ClientTimeout: 1 * time.Second,
		},
	}
}

func (suite *BufconnSuite) SetupBufconn(opts ...grpc.ServerOption) {
	suite.grpcServer = grpc.NewServer(opts...)
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
	ctx, cancel := context.WithTimeout(context.Background(), suite.ClientTimeout)
	defer cancel()

	return grpc.DialContext(
		ctx,
		serverName,
		append(
			opts,
			grpc.WithContextDialer(suite.contextDialer),
			grpc.WithReturnConnectionError(),
		)...,
	)
}
