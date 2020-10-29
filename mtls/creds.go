package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type CredsFiles struct {
	Cert, Key, PeerCACert string
}

type loadedCredsFiles struct {
	cert       tls.Certificate
	peerCAPool *x509.CertPool
}

func loadCredsFiles(credsFiles CredsFiles) (*loadedCredsFiles, error) {
	cert, err := tls.LoadX509KeyPair(credsFiles.Cert, credsFiles.Key)
	if err != nil {
		return nil, err
	}

	peerCACert, err := ioutil.ReadFile(credsFiles.PeerCACert)
	if err != nil {
		return nil, err
	}

	peerCAPool := x509.NewCertPool()
	if ok := peerCAPool.AppendCertsFromPEM(peerCACert); !ok {
		return nil, errors.New("failed to append to peer CA cert pool")
	}

	return &loadedCredsFiles{cert, peerCAPool}, nil
}

func NewServerCreds(serverFiles CredsFiles) (credentials.TransportCredentials, error) {
	loaded, err := loadCredsFiles(serverFiles)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{loaded.cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    loaded.peerCAPool,
	}
	return credentials.NewTLS(config), nil
}

func NewClientCreds(clientFiles CredsFiles) (credentials.TransportCredentials, error) {
	loaded, err := loadCredsFiles(clientFiles)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{loaded.cert},
		RootCAs:      loaded.peerCAPool,
	}
	return credentials.NewTLS(config), nil
}
