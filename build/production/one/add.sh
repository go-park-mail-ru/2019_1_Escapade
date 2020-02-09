#!/bin/sh
echo ""
echo "  -------------------------------------"
echo "  ----Add droplet to docker-machine----"
echo "  -------------------------------------"
echo ""

# Input:
ip=$1 # the IP of the node
path=$2 # path to ssh key
name=$3 # new name of machine

if [ -z "${name}" ]; then
    chmod +x ./../print_error.sh && ./../print_error.sh \
        "$0 <machine ip> <path to id_rsa> <new name of machine>"
    exit 1
fi

yes | docker-machine rm $name
docker-machine create --driver=generic \
    --generic-ip-address=$ip \
    --generic-ssh-user=root \
    --generic-ssh-key=$HOME/$path \
        $name