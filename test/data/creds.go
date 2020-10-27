package testdata

import (
	"google.golang.org/grpc/credentials"

	"github.com/andrejtokarcik/jobworker/mtls"
)

func DefaultServerCredsFiles() mtls.CredsFiles {
	return mtls.CredsFiles{
		Cert:       "server-ca/server1.crt",
		Key:        "server-ca/server1.key",
		PeerCACert: "client-ca.crt",
	}
}

func DefaultClientCredsFiles() mtls.CredsFiles {
	return mtls.CredsFiles{
		Cert:       "client-ca/client1.crt",
		Key:        "client-ca/client1.key",
		PeerCACert: "server-ca.crt",
	}
}

func CredsFilesPaths(rel mtls.CredsFiles) (abs mtls.CredsFiles) {
	abs.Cert = X509Path(rel.Cert)
	abs.Key = X509Path(rel.Key)
	abs.PeerCACert = X509Path(rel.PeerCACert)
	return
}

func DefaultServerCreds() credentials.TransportCredentials {
	creds, err := mtls.NewServerCreds(
		CredsFilesPaths(DefaultServerCredsFiles()),
	)
	if err != nil {
		panic(err)
	}
	return creds
}

func DefaultClientCreds() credentials.TransportCredentials {
	creds, err := mtls.NewClientCreds(
		CredsFilesPaths(DefaultClientCredsFiles()),
	)
	if err != nil {
		panic(err)
	}
	return creds
}
