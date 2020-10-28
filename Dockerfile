# syntax = docker/dockerfile:1.1-experimental

ARG PRODUCTS_OF
ARG GO_VERSION
ARG MOCKERY_VERSION

FROM golang:${GO_VERSION}-buster AS base
RUN apt-get update && apt-get install -y --no-install-recommends unzip
WORKDIR /src
ENV OUTROOT /out
RUN mkdir $OUTROOT
ENV MAKELOCAL make -f Makefile.local
COPY go.mod go.sum .
RUN go mod download

FROM base AS build
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    $MAKELOCAL build

FROM base AS test
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    $MAKELOCAL test

FROM base AS protoc-base
ARG PROTOC_VERSION
ARG GOGO_PROTO_VERSION
ENV PROTOC_ZIP protoc-${PROTOC_VERSION}-linux-x86_64.zip
RUN curl -sSLO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP} && \
    mkdir -p /opt/protoc && unzip -d /opt/protoc $PROTOC_ZIP && \
    install -Ds /opt/protoc/bin/protoc /usr/bin/protoc
RUN mkdir -p ${GOPATH}/src/github.com/gogo/protobuf && \
    curl -sSL https://api.github.com/repos/gogo/protobuf/tarball/${GOGO_PROTO_VERSION} | tar xz --strip 1 -C ${GOPATH}/src/github.com/gogo/protobuf && \
    cd ${GOPATH}/src/github.com/gogo/protobuf && \
    go build -ldflags '-w -s' -o /opt/gogo-protobuf-out/protoc-gen-gogofaster ./protoc-gen-gogofaster && \
    install -Ds /opt/gogo-protobuf-out/protoc-gen-gogofaster /usr/bin/protoc-gen-gogofaster && \
    install -D $(find ./protobuf/google/protobuf -name '*.proto') -t /usr/include/proto/google/protobuf && \
    install -D ./gogoproto/gogo.proto /usr/include/proto/github.com/gogo/protobuf/gogoproto/gogo.proto
ENV PROTOC_INCLUDE /usr/include/proto

FROM vektra/mockery:v${MOCKERY_VERSION} AS mockery-base

FROM protoc-base AS gen-base
COPY --from=mockery-base /usr/local/bin/mockery /usr/bin/mockery

FROM gen-base AS gen
RUN --mount=target=. \
    $MAKELOCAL gen

FROM gen-base AS check-gen
RUN --mount=target=.,readwrite \
    --mount=type=cache,target=/root/.cache/go-build \
    $MAKELOCAL check-gen

FROM $PRODUCTS_OF AS create-products
FROM scratch AS copy-products
COPY --from=create-products /out/ /
