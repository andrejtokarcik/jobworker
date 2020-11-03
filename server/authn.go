package server

import (
	"context"
	"crypto/x509/pkix"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type clientSubjectKey struct{}

// ClientSubject identifies a client during a communication with the server.
type ClientSubject = pkix.Name

func AttachClientSubject(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no peer found")
	}

	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unexpected peer transport credentials")
	}

	if len(tlsAuth.State.PeerCertificates) == 0 {
		return nil, status.Error(codes.Unauthenticated, "could not verify peer certificate")
	}

	newCtx := context.WithValue(ctx, clientSubjectKey{}, tlsAuth.State.PeerCertificates[0].Subject)
	return handler(newCtx, req)
}

func getClientSubject(ctx context.Context) (ClientSubject, error) {
	value := ctx.Value(clientSubjectKey{})
	if value == nil {
		return ClientSubject{}, status.Errorf(codes.Internal, "cannot obtain client subject")
	}

	clientSubject, ok := value.(ClientSubject)
	if !ok {
		return ClientSubject{}, status.Errorf(codes.Internal, "cannot determine client subject")
	}

	return clientSubject, nil
}
