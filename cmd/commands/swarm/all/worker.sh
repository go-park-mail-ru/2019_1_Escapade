#!/bin/sh
echo ""
echo "  ----------------------------"
echo "  -----Make nodes workers-----"
echo "  ----------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x prepare.sh && ./prepare.sh

worker() {
    swarm_token="SWMTKN-1-0njltu4byk10q82oxkp4097tf7wvaog5h1q1od6addrqqay91s-2ott62721zc6rs6ma8fs3jh8x"
    script="./one/worker.sh"
    chmod +x $script && \
    $script $1 $swarm_token $2
}

# $1 - адрес менеджера
# $2 - $5 - адреса будущих воркеров
worker $2 $1 && \
worker $3 $1 && \
worker $4 $1 && \
worker $5 $1

