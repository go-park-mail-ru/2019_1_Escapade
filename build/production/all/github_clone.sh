#!/bin/sh
echo ""
echo "  -----------------------------------"
echo "  -----$1 Firewall set-----"
echo "  -----------------------------------"
echo ""
trap 'echo " stop" ' INT TERM

echo "  Prepare /one/github_clone.sh"
chmod +x ./one/github_clone.sh

while [ $1 ]
    do
    ./one/github_clone.sh $1
    shift 1
done
