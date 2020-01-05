#!/bin/bash
echo ""
echo "  -----------------------------"
echo "  ---Make $1 manager---"
echo "  -----------------------------"
echo ""

#$1 - the IP of this node
ssh root@$1 "
docker swarm leave --force
docker swarm init --advertise-addr $1:2377
docker network rm backend-overlay
docker network create -d overlay --subnet 10.10.9.0/24 --attachable backend-overlay
docker node ls
 "

