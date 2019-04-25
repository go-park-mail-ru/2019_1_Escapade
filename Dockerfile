FROM golang:alpine
WORKDIR /escapade
COPY . .
RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN go build -o ../api main.go