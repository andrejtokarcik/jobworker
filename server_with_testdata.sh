#!/bin/bash

export GRPC_GO_LOG_SEVERITY_LEVEL=${LOG_LEVEL:-warning}

CERT_DIR=${CERT_DIR:-./mtls/testdata}

./bin/server \
    --server-cert $CERT_DIR/server-ca/server1.crt \
    --server-key $CERT_DIR/server-ca/server1.key \
    --client-ca-cert $CERT_DIR/client-ca.crt $@
