package server

import (
	"context"
	"crypto/x509/pkix"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RPCCallerFilter is a function that returns true if the given client
// (identified by a certificate subject) should be allowed to call the
// given RPC method, false otherwise.
type RPCCallerFilter func(pkix.Name, *grpc.UnaryServerInfo) bool

// ApplyRPCCallerFilter returns a unary server interceptor that may reject
// a client with "permission denied" based on the given RPCCallerFilter.
func ApplyRPCCallerFilter(filter RPCCallerFilter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		value := ctx.Value(clientSubjectKey{})
		if value == nil {
			return nil, status.Errorf(codes.Internal, "cannot determine client subject")
		}

		if !filter(value.(pkix.Name), info) {
			return nil, status.Errorf(codes.PermissionDenied, "not allowed to call this RPC method")
		}
		return handler(ctx, req)
	}
}
