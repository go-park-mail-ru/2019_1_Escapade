#!/bin/sh
echo ""
echo "  ------------------------"
echo "  -----Make manager ------"
echo "  ------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x droplet_manager.sh && ./droplet_manager.sh 1

#$1 - the IP of this node

ssh root@$1 "
docker swarm init --advertise-addr $1 
docker network rm backend-overlay
docker network create -d overlay --attachable backend-overlay
 "
