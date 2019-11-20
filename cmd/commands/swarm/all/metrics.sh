#!/bin/sh
echo ""
echo "  ----------------------------"
echo "  -----Nodes metrics fix-----"
echo "  ----------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x prepare.sh && ./prepare.sh

metrics() {
    script="./one/metrics.sh"
    chmod +x $script && \
    $script $1
}

# $1 - $5 - адреса серверов
metrics $1 && \
metrics $2 && \
metrics $3 && \
metrics $4 && \
metrics $5 

 
