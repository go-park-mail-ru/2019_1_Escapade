# syntax = docker/dockerfile:experimental
FROM golang:alpine as build
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR

COPY go.mod .
COPY go.sum .

# RUN apk add --update git
# RUN apk add --update bash && rm -rf /var/cache/apk/*
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
COPY /docs /2019_1_Escapade/docs

COPY /cmd/services/auth .
RUN --mount=type=cache,target=/root/.cache/go-build \ 
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /go/bin/auth .

FROM scratch
# FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/auth /
ENTRYPOINT ["./auth"]  