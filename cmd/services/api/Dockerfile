# syntax = docker/dockerfile:experimental
FROM golang:alpine as build
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR

# introducing DOCKER_BUILDKIT=1
#    first  build(no cache)   - 225s
#    second build(with cache) - 64s

COPY go.mod .
COPY go.sum .

RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN --mount=type=cache,target=/root/.cache/go-build \ 
    --mount=type=cache,target=/go/pkg/mod \ 
    go mod download

RUN mkdir -p /2019_1_Escapade/internal/pkg
RUN mkdir -p /2019_1_Escapade/internal/services/api
RUN mkdir /2019_1_Escapade/docs

COPY /internal/pkg /2019_1_Escapade/internal/pkg
COPY /internal/services/api /2019_1_Escapade/internal/services/api
COPY /docs/api /2019_1_Escapade/docs/api

COPY /cmd/services/api .
RUN --mount=type=cache,target=/root/.cache/go-build \ 
    --mount=type=cache,target=/go/pkg/mod \
    go build -o bin/api main.go

FROM scratch
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR
COPY --from=build $SRC_DIR/bin/api $SRC_DIR/bin/api
ENTRYPOINT ["./bin/api"]  
