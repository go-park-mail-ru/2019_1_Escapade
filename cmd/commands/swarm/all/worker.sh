#!/bin/sh
echo ""
echo "  ----------------------------"
echo "  -----Make nodes workers-----"
echo "  ----------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x prepare.sh && ./prepare.sh

worker() {
    swarm_token="SWMTKN-1-3w0ngu7m0o1sje8pkqfzw7mft0v6s4hqs53m97e6jddslqhmkn-48iywq96luawzkybx4ksoqa8t"
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

