# set consul version
ENV CONSUL_VERSION 1.6.1

# base image
FROM progrium/consul:latest

# create a new directory
RUN mkdir /consul
RUN mkdir /consul/config

# add the config file
COPY . /consul/config