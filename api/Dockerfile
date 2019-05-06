FROM ubuntu:18.04
WORKDIR /escapade
COPY go.mod .
COPY go.sum .
RUN apt-get -y update
RUN apt install -y gcc git wget

RUN wget https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz
RUN tar -xvf go1.12.4.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

RUN go mod download
RUN go build -o bin/api api/main.go
COPY . .
CMD go test ./...