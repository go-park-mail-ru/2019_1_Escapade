#!/bin/sh
echo ""
echo "  ----------------------------"
echo "  -----Nodes Firewall set-----"
echo "  ----------------------------"
echo ""

echo "  Prepare /one/firewall.sh"
chmod +x ./one/firewall.sh

while [ $1 ]
    do
    ./one/firewall.sh $1
    shift 1
done
