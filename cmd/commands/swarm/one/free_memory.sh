#!/bin/sh
echo ""
echo "  -------------------------------"
echo "  ---------Free memory-----------"
echo "  -------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x free_memory.sh && ./free_memory.sh

#$1 - the IP of the node
#$2 - path to ssh key
#$3 - new name of machine

echo "before"
free -m

sync; echo 2 > /proc/sys/vm/drop_caches
sync; echo 3 > /proc/sys/vm/drop_caches

echo "after"
free -m
