#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------- Prepare env -------"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#source ./make_env.sh

IP=`ifconfig eth0 | awk 'BEGIN{FS="\n"; RS=""} {print $2}' |  awk ' {print $2}'`
export IP
echo "IP="$IP
















