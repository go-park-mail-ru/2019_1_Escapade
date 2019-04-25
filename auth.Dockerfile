FROM golang:alpine
WORKDIR /escapade-auth
COPY . .
RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN go build -o bin/auth auth/auth.go