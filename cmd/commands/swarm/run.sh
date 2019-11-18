#!/bin/sh
echo ""
echo "  --------------------------"
echo "  -Create images and deploy-"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x run.sh && ./run.sh
chmod +x images.sh
cd ../../..
#sudo service docker restart
#sudo docker-compose build
docker run -d -p 5000:5000 --restart=always --name registry registry:2
#sudo docker-compose push
./images.sh
sudo docker stack deploy -c docker-swarm.yaml app
docker container stop registry
