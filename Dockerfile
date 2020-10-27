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

FROM vektra/mockery:v${MOCKERY_VERSION} AS generate-base

FROM base AS generate
RUN --mount=target=. \
    --mount=from=generate-base,src=/usr/local/bin/mockery,target=/usr/bin/mockery \
    --mount=type=cache,target=/root/.cache/go-build \
    $MAKELOCAL generate

FROM base AS protoc-base
ARG PROTOC_VERSION
ARG GOGO_PROTO_VERSION
ENV PROTOC_ZIP protoc-${PROTOC_VERSION}-linux-x86_64.zip
ENV PROTOC_URL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}
ENV PROTOC_DIR /opt/protoc
RUN curl -LO $PROTOC_URL
RUN mkdir $PROTOC_DIR && unzip -d $PROTOC_DIR $PROTOC_ZIP
RUN go get github.com/gogo/protobuf/protoc-gen-gogo@${GOGO_PROTO_VERSION}
ENV PROTOC ${PROTOC_DIR}/bin/protoc

FROM protoc-base AS protoc
RUN --mount=target=. \
    $MAKELOCAL protoc

FROM protoc-base AS check-gen
RUN --mount=target=.,readwrite \
    --mount=from=generate-base,src=/usr/local/bin/mockery,target=/usr/bin/mockery \
    --mount=type=cache,target=/root/.cache/go-build \
    $MAKELOCAL check-gen

FROM $PRODUCTS_OF AS gen-products
FROM scratch AS copy-products
COPY --from=gen-products /out/ /
