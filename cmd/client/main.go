package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/andrejtokarcik/jobworker/client"
	"github.com/andrejtokarcik/jobworker/mtls"
)

var (
	serverAddress string
	connTimeout   time.Duration
	credsFiles    mtls.CredsFiles
)

func init() {
	flag.StringVar(&serverAddress, "server", "127.0.0.1:50051", "Address of the server to connect to")
	flag.DurationVar(&connTimeout, "timeout", 5*time.Second, "Connection timeout")

	flag.StringVar(&credsFiles.Cert, "client-cert", "client.crt", "Certificate file to use for the client")
	flag.StringVar(&credsFiles.Key, "client-key", "client.key", "Private key file to use for the client")
	flag.StringVar(&credsFiles.PeerCACert, "server-ca-cert", "server-ca.crt", "Certificate file of the CA to authenticate the server")
}

func main() {
	flag.Parse()

	creds, err := mtls.NewClientCreds(credsFiles)
	if err != nil {
		log.Fatal("Failed to load mTLS credentials: ", err)
	}

	conn, err := client.DialContextWithTimeout(
		context.Background(),
		connTimeout,
		serverAddress,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatal("Failed to dial server: ", err)
	}
	defer conn.Close()

	log.Print("Successfully connected to server at ", serverAddress)
}
