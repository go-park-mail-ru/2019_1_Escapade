#!/bin/sh
echo ""
echo "  ----------------------------"
echo "  ----Docker swarm prepare----"
echo "  ----------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x prepare.sh && ./prepare.sh

# scripts from a folder 'swarm/manager' can be called on manager node
# scripts from a folder 'swarm/one' can be called on any internal node
# scripts from a folder 'swarm/all' can be called on external machine
# scripts from a folder 'swarm' can be called on external machine

echo "  0. Prepare other .sh"
chmod +x ./all/firewall.sh && \
chmod +x ./all/diagnostics.sh && \
chmod +x ./all/metrics.sh && \
chmod +x ./one/manager.sh && \
chmod +x ./all/worker.sh  

# адреса серверов
addr1="68.183.48.80"
addr2="167.71.247.116"
addr3="167.172.21.178"
addr4="167.172.21.125"
addr5="206.81.9.205"

./all/firewall.sh $addr1 $addr2 $addr3 $addr4 $addr5 && \
./all/diagnostics.sh $addr1 $addr2 $addr3 $addr4 $addr5 && \
./all/metrics.sh $addr1 $addr2 $addr3 $addr4 $addr5 && \
./one/manager.sh && \ #dont forget to set manager token to worker.sh 
./all/worker.sh $addr1 $addr2 $addr3 $addr4 $addr5