#!/bin/sh
echo ""
echo "  --------------------------"
echo "  -Create images and deploy-"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x run.sh && ./run.sh

cd ../../..
sudo service docker restart
sudo docker-compose build
sudo docker-compose push
sudo docker stack deploy -c docker-compose.yaml app