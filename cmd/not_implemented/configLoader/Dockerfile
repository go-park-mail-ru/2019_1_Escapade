# docker build -t api .

FROM golang:alpine
WORKDIR /escapade
COPY go.mod .
COPY go.sum .
RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*

RUN go mod download

# download dependencies
RUN apk --no-cache add \
      bash \
      ca-certificates \
      wget

ENV VAULT_VERSION 0.10.3

# download and set up vault
RUN wget --quiet --output-document=/tmp/vault.zip https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip 
RUN    unzip /tmp/vault.zip -d /escapade 
RUN rm -f /tmp/vault.zip && \
    chmod +x /escapade/vault

COPY cmd/services/configLoader .

RUN go build -o bin/loader main.go