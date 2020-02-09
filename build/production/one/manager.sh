#!/bin/bash
echo ""
echo "  -----------------------------"
echo "  ---Make $1 manager---"
echo "  -----------------------------"
echo ""

ip=$1   # the IP of this node

if [ -z "${ip}" ]; then
    chmod +x ./../print_error.sh && ./../print_error.sh \
        "$0 <ip of the node>"
    exit 1
fi

ssh root@$ip "
docker swarm leave --force
docker swarm init --advertise-addr $1:2377
docker network rm backend-overlay
docker network create -d overlay --subnet 10.10.9.0/24 --attachable backend-overlay
docker node ls
 "

