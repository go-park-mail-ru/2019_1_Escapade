# docker build -t api .

FROM golang:alpine
WORKDIR /escapade
COPY go.mod .
COPY go.sum .
RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN go mod download

RUN mkdir -p /internal/pkg
COPY /internal/pkg /internal/pkg

RUN mkdir -p /internal/services/api
COPY /internal/services/api /internal/services/api

COPY /cmd/services/api .
RUN go build -o bin/api main.go
ENTRYPOINT ["./bin/api"]  