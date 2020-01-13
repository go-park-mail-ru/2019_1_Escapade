#!/bin/sh
#chmod +x node_drain.sh && ./node_drain.sh 1

node=$1   # the name of this node

if [ -z "${node}" ]; then
    chmod +x ./../print_error.sh && ./../print_error.sh \
        "make drain node=<name>"
    exit 1
fi

docker node update --availability Drain ${node}