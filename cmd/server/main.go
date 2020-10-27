package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/andrejtokarcik/jobworker/mtls"
)

var (
	grpcPort   int
	credsFiles mtls.CredsFiles
)

func init() {
	flag.IntVar(&grpcPort, "grpc-port", 50051, "Port to expose the gRPC server on")

	flag.StringVar(&credsFiles.Cert, "server-cert", "server.crt", "Certificate file to use for the server")
	flag.StringVar(&credsFiles.Key, "server-key", "server.key", "Private key file to use for the server")
	flag.StringVar(&credsFiles.PeerCACert, "client-ca-cert", "client-ca.crt", "Certificate file of the CA to authenticate the clients")
}

func main() {
	flag.Parse()

	creds, err := mtls.NewServerCreds(credsFiles)
	if err != nil {
		log.Fatal("Failed to load mTLS credentials: ", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	log.Print("Starting gRPC server at ", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
