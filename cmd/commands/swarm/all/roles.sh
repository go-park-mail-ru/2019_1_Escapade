#!/bin/sh
echo ""
echo "  ------------------------------"
echo "  ---Make manager and workers---"
echo "  ------------------------------"
echo ""

chmod +x ./one/manager.sh && chmod +x ./one/worker.sh

manager=$1
./one/manager.sh $manager
shift 1

token=`ssh root@$manager "docker swarm join-token worker -q"`

while [ $1 ]
    do
    ./one/worker.sh $1 $token $manager
    shift 1
done
