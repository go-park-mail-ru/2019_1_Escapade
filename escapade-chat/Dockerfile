FROM golang:alpine
WORKDIR /escapade
COPY go.mod .
COPY go.sum .
RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN go mod download
COPY . .
RUN go build -o bin/chat main.go