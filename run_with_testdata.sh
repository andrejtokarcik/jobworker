#!/bin/bash

trap 'kill $SERVER_PID' TERM INT EXIT
export GRPC_GO_LOG_SEVERITY_LEVEL=warning

CERT_DIR=${CERT_DIR:-./mtls/testdata}

./bin/server \
    -server-cert $CERT_DIR/server-ca/server1.crt \
    -server-key $CERT_DIR/server-ca/server1.key \
    -client-ca-cert $CERT_DIR/client-ca.crt &
SERVER_PID=$!

sleep 1

./bin/client \
    -client-cert $CERT_DIR/client-ca/client1.crt \
    -client-key $CERT_DIR/client-ca/client1.key \
    -server-ca-cert $CERT_DIR/server-ca.crt $@
