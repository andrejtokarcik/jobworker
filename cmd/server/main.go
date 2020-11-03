package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/andrejtokarcik/jobworker/mtls"
	"github.com/andrejtokarcik/jobworker/server"
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
	grpcServer := server.New(filterRPCCallers, grpc.Creds(creds))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	log.Print("Starting gRPC server at ", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}

// XXX The rules encoded in filterRPCCallers would be ordinarily inferred
// from a configuration file, e.g. by means of https://github.com/spf13/viper
func filterRPCCallers(client server.ClientSubject, rpcInfo *grpc.UnaryServerInfo) bool {
	switch rpcInfo.FullMethod {
	case jobWorkerMethod("StartJob"):
		return client.CommonName != "client2"
	case jobWorkerMethod("StopJob"):
		return client.CommonName != "client2"
	case jobWorkerMethod("GetJob"):
		return true
	default:
		return false
	}
}

func jobWorkerMethod(rpc string) string {
	return fmt.Sprintf("/jobworker.JobWorker/%s", rpc)
}
