#!/bin/sh
echo ""
echo "  ------------------------------"
echo "  --Make $1 worker--"
echo "  ------------------------------"
echo ""

#$1 - the IP of this node
#$2 - the manager token
#$3 - the IP of the Swarm manager

ssh root@$1 "
docker swarm leave --force
docker swarm join --advertise-addr $1:2377 --token $2 $3:2377
"