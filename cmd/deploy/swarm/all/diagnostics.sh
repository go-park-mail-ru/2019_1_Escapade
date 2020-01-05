#!/bin/sh
echo ""
echo "  ------------------------------"
echo "  ----Nodes check ports open----"
echo "  ------------------------------"
echo ""

echo "  Prepare /one/diagnostics.sh"
chmod +x ./one/diagnostics.sh

manager=$1
shift 1
while [ $1 ]
    do
    ./one/diagnostics.sh $manager $1
    ./one/diagnostics.sh $1 $manager
    shift 1
done
