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
    CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /go/bin/game .

FROM scratch
# FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/game /
ENTRYPOINT ["./game"]  