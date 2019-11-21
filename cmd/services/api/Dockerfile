# docker build -t api .

FROM golang:alpine as build
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR

COPY go.mod .
COPY go.sum .

RUN apk add --update git
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN go mod download

#RUN mkdir src && cd src && \
#wget http://mirror.linux-ia64.org/gnu/gcc/releases/gcc-8.3.0/gcc-8.3.0.tar.gz && \
#tar -xvf gcc-8.3.0.tar.gz && mkdir bld && cd bld && \
#../src/configure --target=powerpc64le-linux --host=powerpc64le-linux --build=powerpc64le-linux --enable-threads=posix --enable-shared --enable-__cxa_atexit --enable-languages=c,c++,go --enable-secureplt --enable-checking=yes --with-long-double-128 --enable-decimal-float --disable-bootstrap --disable-alsa --disable-multilib –prefix=/usr/local && \
#make && \ make install

RUN mkdir -p /2019_1_Escapade/internal/pkg
COPY /internal/pkg /2019_1_Escapade/internal/pkg

RUN mkdir -p /2019_1_Escapade/internal/services/api
COPY /internal/services/api /2019_1_Escapade/internal/services/api

RUN mkdir -p /2019_1_Escapade/docs
COPY /docs/api /2019_1_Escapade/docs/api

COPY /cmd/services/api .
RUN go build -o bin/api main.go

FROM scratch
ENV SRC_DIR=/2019_1_Escapade/
WORKDIR $SRC_DIR
COPY --from=build $SRC_DIR/bin/api $SRC_DIR/bin/api
ENTRYPOINT ["./bin/api"]  
