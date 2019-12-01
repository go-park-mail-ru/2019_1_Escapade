#!/bin/sh
echo ""
echo "  -------------------------------"
echo "  -Add droplet to docker-machine-"
echo "  -------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x droplet_add.sh && ./droplet_add.sh

#$1 - the IP of the node
#$2 - path to ssh key
#$3 - new name of machine

yes | docker-machine rm $3
docker-machine create --driver=generic \
    --generic-ip-address=$1 \
    --generic-ssh-user=root \
    --generic-ssh-key=$2 \
        $3