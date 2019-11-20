#!/bin/sh
echo ""
echo "  ----------------------------"
echo "  ----Docker swarm prepare----"
echo "  ----------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x prepare.sh && ./prepare.sh

echo "  0. Prepare other .sh"
chmod +x /all/firewall.sh && \
chmod +x /all/diagnostics.sh && \
chmod +x /all/metrics.sh && \
chmod +x /one/manager.sh && \
chmod +x /all/worker.sh  

# адреса серверов
addr1="142.93.79.194"
addr2="167.71.247.116"
addr3="167.172.21.178"
addr4="167.172.21.125"
addr5="206.81.9.205"

./all/firewall.sh $addr1 $addr2 $addr3 $addr4 $addr5 && \
./all/diagnostics.sh $addr1 $addr2 $addr3 $addr4 $addr5 && \
./all/metrics.sh $addr1 $addr2 $addr3 $addr4 $addr5 && \
./one/manager.sh && \ #dont forget to set manager token to worker.sh 
./all/worker.sh $addr1 $addr2 $addr3 $addr4 $addr5