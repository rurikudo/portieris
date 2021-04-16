FROM golang:1.16.3 as golang


WORKDIR /workspace

COPY . ./

COPY go.mod go.mod
COPY go.sum go.sum


RUN go mod tidy
RUN go mod download

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o /go/bin/main \
    -ldflags '-s -w'


# FROM busybox
FROM ubuntu
ENV DOCKER_CONFIG /tmp
# ENV COSIGN_EXPERIMENTAL "1"
COPY --from=golang /go/bin/main /app/main
RUN apt-get update
# RUN apt-get install -y curl jq openssl
RUN apt-get install -y openssl curl
ENTRYPOINT ["/app/main"]




# # FROM registry.access.redhat.com/ubi7/ubi:latest
# FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

# # RUN yum install -y tar

# RUN mkdir -p /cosign-app
# COPY . /cosign-app

# RUN chgrp -R 0 /cosign-app && chmod -R g=u /cosign-app

# WORKDIR /cosign-app

# ENTRYPOINT ["/cosign-app"]