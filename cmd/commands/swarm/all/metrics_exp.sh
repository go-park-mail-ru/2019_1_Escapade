#!/bin/sh
echo ""
echo "  ----------------------------"
echo "  -----Nodes metrics fix-----"
echo "  ----------------------------"
echo ""

echo "  Prepare /one/metrics.sh"
chmod +x ./one/metrics.sh

shift 1
while [ $1 ]
    do
    ./one/metrics.sh $1
    shift 1
done
