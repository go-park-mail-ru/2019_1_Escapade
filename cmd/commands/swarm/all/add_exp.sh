#!/bin/sh
echo ""
echo "  -------------------------------"
echo "  --Add nodes to docker-machine--"
echo "  -------------------------------"
echo ""

echo "  Prepare /one/add.sh"
chmod +x ./one/add.sh

while [ $1 ]
    do
    ./one/add.sh $1 $2 $3
    shift 3
done