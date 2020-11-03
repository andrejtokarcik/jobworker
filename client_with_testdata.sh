#!/bin/sh

export GRPC_GO_LOG_SEVERITY_LEVEL=${LOG_LEVEL:-warning}

CERT_DIR=${CERT_DIR:-./mtls/testdata}
CLIENT=${CLIENT:-client1}

./bin/client \
    --client-cert $CERT_DIR/client-ca/$CLIENT.crt \
    --client-key $CERT_DIR/client-ca/$CLIENT.key \
    --server-ca-cert $CERT_DIR/server-ca.crt $@
