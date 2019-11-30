FROM golang:alpine as build
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR

COPY go.mod .
COPY go.sum .

RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN go mod download

RUN mkdir -p /2019_1_Escapade/internal/pkg
COPY /internal/pkg /2019_1_Escapade/internal/pkg

RUN mkdir -p /2019_1_Escapade/internal/services/api
COPY /internal/services/api /2019_1_Escapade/internal/services/api

RUN mkdir -p /2019_1_Escapade/internal/services/auth
COPY /internal/services/auth /2019_1_Escapade/internal/services/auth

RUN mkdir -p /2019_1_Escapade/internal/services/ery
COPY /internal/services/ery /2019_1_Escapade/internal/services/ery

RUN mkdir -p /2019_1_Escapade/docs
COPY /docs/auth /2019_1_Escapade/docs/auth

COPY /cmd/services/auth .
RUN go build -o bin/auth main.go

FROM scratch
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR
COPY --from=build $SRC_DIR/bin/auth $SRC_DIR/bin/auth
ENTRYPOINT ["./bin/auth"]  