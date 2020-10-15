# syntax = docker/dockerfile:1.1-experimental

ARG PRODUCTS_OF
ARG GO_VERSION

FROM golang:${GO_VERSION}-buster AS base
RUN apt-get update && apt-get install -y --no-install-recommends unzip
WORKDIR /src
ENV OUTROOT /out
RUN mkdir $OUTROOT
ENV MAKE make -f Makefile.local
COPY go.mod go.sum .
RUN go mod download

FROM base AS build
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    $MAKE build

FROM base AS test
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    $MAKE test

FROM base AS check-gen
RUN --mount=target=.,readwrite \
    --mount=type=cache,target=/root/.cache/go-build \
    $MAKE check-gen

FROM $PRODUCTS_OF AS gen-products
FROM scratch AS copy-products
COPY --from=gen-products /out/ /
