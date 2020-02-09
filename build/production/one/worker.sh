#!/bin/sh
echo ""
echo "  ------------------------------"
echo "  --Make $1 worker--"
echo "  ------------------------------"
echo ""

ip=$1        # the IP of this node
token=$2     # the manager token
manager=$3   # the IP of the Swarm manager

if [ -z "${manager}" ]; then
    chmod +x ./../print_error.sh && ./../print_error.sh \
        "$0 <ip of the node> <the manager token> <the IP of the Swarm manager>"
    exit 1
fi

ssh root@$ip "
docker swarm leave --force
docker swarm join --advertise-addr $ip --token $token $manager:2377
"