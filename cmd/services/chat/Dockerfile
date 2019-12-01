# syntax = docker/dockerfile:experimental
FROM golang:alpine as build
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR

# introducing DOCKER_BUILDKIT=1
#    with cache - 34s

COPY go.mod .
COPY go.sum .

RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN --mount=type=cache,target=/root/.cache/go-build \ 
    --mount=type=cache,target=/go/pkg/mod \ 
    go mod download

RUN mkdir -p /2019_1_Escapade/internal/pkg
RUN mkdir -p /2019_1_Escapade/internal/services/chat

COPY /internal/pkg /2019_1_Escapade/internal/pkg
COPY /internal/services/chat /2019_1_Escapade/internal/services/chat

COPY /cmd/services/chat .
RUN --mount=type=cache,target=/root/.cache/go-build \ 
    --mount=type=cache,target=/go/pkg/mod \
    go build -o bin/chat main.go

FROM scratch
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR
COPY --from=build $SRC_DIR/bin/chat $SRC_DIR/bin/chat
ENTRYPOINT ["./bin/chat"]  