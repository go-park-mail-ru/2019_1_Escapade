#!/bin/sh
echo ""
echo "  --------------------------------"
echo "  ---Updating docker containers---"
echo "  --------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#sudo chmod +x docker_compose.sh && ./docker_compose.sh

sudo docker-compose stop
sudo docker-compose rm -f
sudo docker-compose pull   
sudo docker-compose up -d