# syntax = docker/dockerfile:experimental
FROM golang:alpine as build
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR

# introducing DOCKER_BUILDKIT=1
# before: 
#    first  build(no cache)   - 369s
#    second build(with cache) - 220s
# after :
#    first  build(no cache)   - 278s
#    second build(with cache) - 112s

COPY go.mod .
COPY go.sum .

RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN --mount=type=cache,target=/root/.cache/go-build \ 
    --mount=type=cache,target=/go/pkg/mod \ 
    go mod download

RUN mkdir /2019_1_Escapade/internal
RUN mkdir /2019_1_Escapade/internal/pkg
RUN mkdir /2019_1_Escapade/internal/services
RUN mkdir /2019_1_Escapade/internal/services/api
RUN mkdir /2019_1_Escapade/internal/services/auth
RUN mkdir /2019_1_Escapade/internal/services/ery
RUN mkdir /2019_1_Escapade/docs

COPY /internal/pkg /2019_1_Escapade/internal/pkg
COPY /internal/services/api /2019_1_Escapade/internal/services/api
COPY /internal/services/auth /2019_1_Escapade/internal/services/auth
COPY /internal/services/ery /2019_1_Escapade/internal/services/ery
COPY /docs/auth /2019_1_Escapade/docs/auth

COPY /cmd/services/auth .
RUN --mount=type=cache,target=/root/.cache/go-build \ 
    --mount=type=cache,target=/go/pkg/mod \
    go build -o bin/auth -i main.go

FROM scratch
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR
COPY --from=build $SRC_DIR/bin/auth $SRC_DIR/bin/auth
ENTRYPOINT ["./bin/auth"]  