#!/bin/sh
echo ""
echo "  ------------------------------"
echo "  --Make $1 worker--"
echo "  ------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x droplet_worker.sh && ./droplet_worker.sh 1 2 3

#$1 - the IP of this node
#$2 - the token you received when called ./droplet_manager.sh 
#$3 - the IP of the Swarm manager

ssh root@$1 "
docker swarm leave
docker swarm join --advertise-addr $1:2377 --token $2 $3:2377
"