package main

import (
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/andrejtokarcik/jobworker/client"
	"github.com/andrejtokarcik/jobworker/mtls"
)

var (
	serverAddress  string
	timeoutSeconds int
	credsFiles     mtls.CredsFiles
)

func init() {
	flag.StringVar(&serverAddress, "server", "0.0.0.0:50051", "Address of the server to connect to")
	flag.IntVar(&timeoutSeconds, "timeout", 5, "Connection timeout in seconds")

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

	conn, err := client.Dial(
		serverAddress,
		time.Duration(timeoutSeconds)*time.Second,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatal("Failed to dial server: ", err)
	}
	defer conn.Close()

	log.Print("Successfully connected to server at ", serverAddress)
}
