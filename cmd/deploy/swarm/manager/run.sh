#!/bin/sh
echo ""
echo "  ----------------------------------------------------"
echo "  ----Create images, push them and deploy as stack----"
echo "  ----------------------------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x run.sh && ./run.sh
chmod +x install.sh
chmod +x swap.sh
chmod +x build.sh
chmod +x push.sh

# cd ../../..
#sudo service docker restart
#./swap.sh
#./install.sh
#./build.sh
#./push.sh
#docker-compose build && ./images.sh
#docker run -d -p 5000:5000 --restart=always --name registry registry:2
#sudo docker-compose push
#./images.sh
sudo docker stack deploy -c docker-swarm.yaml app
