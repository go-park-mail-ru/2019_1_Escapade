# syntax = docker/dockerfile:experimental
FROM golang:alpine as build
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR

# introducing DOCKER_BUILDKIT=1
#    without cache - 60s
#    with cache    - 41s

COPY go.mod .
COPY go.sum .

RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN --mount=type=cache,target=/root/.cache/go-build \ 
    --mount=type=cache,target=/go/pkg/mod \ 
    go mod download

RUN mkdir -p /2019_1_Escapade/internal/pkg
RUN mkdir -p /2019_1_Escapade/internal/services/game
RUN mkdir -p /2019_1_Escapade/internal/services/chat
RUN mkdir -p /2019_1_Escapade/internal/services/api

COPY /internal/pkg /2019_1_Escapade/internal/pkg
COPY /internal/services/game /2019_1_Escapade/internal/services/game
COPY /internal/services/chat /2019_1_Escapade/internal/services/chat
COPY /internal/services/api /2019_1_Escapade/internal/services/api


COPY /cmd/services/game .
RUN --mount=type=cache,target=/root/.cache/go-build \ 
    --mount=type=cache,target=/go/pkg/mod \
    go build -o bin/game main.go

FROM scratch
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR
COPY --from=build $SRC_DIR/bin/game $SRC_DIR/bin/game
ENTRYPOINT ["./bin/game"]  