# docker build -t game .

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

RUN mkdir -p /2019_1_Escapade/internal/services/game
COPY /internal/services/game /2019_1_Escapade/internal/services/game

RUN mkdir -p /2019_1_Escapade/internal/services/chat
COPY /internal/services/chat /2019_1_Escapade/internal/services/chat

RUN mkdir -p /2019_1_Escapade/internal/services/api
COPY /internal/services/api /2019_1_Escapade/internal/services/api


COPY /cmd/services/game .
RUN go build -o bin/game main.go

FROM scratch
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR
COPY --from=build $SRC_DIR/bin/game $SRC_DIR/bin/game
ENTRYPOINT ["./bin/game"]  