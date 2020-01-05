#!/bin/sh
echo ""
echo "  ----------------------"
echo "  ----Clear dangling----"
echo "  ----------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x clear.sh && ./clear.sh

docker image ls
docker rm  $(docker ps -q -a)
docker rmi $(docker images -f "dangling=true" -q)
docker image ls