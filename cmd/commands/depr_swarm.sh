#!/bin/sh
echo ""
echo "  ---------------------------"
echo "  ---Launch consul servers---"
echo "  ---------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x swarm.sh && ./swarm.sh

#docker-machine create nb-consul --driver virtualbox

#cat ~/bin/dm-env
#eval `docker-machine env $2 $1`
# avoid typing too much
#alias dm=docker-machine

#docker run -d --restart always -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302/udp -p 8302:8302 -p 8400:8400 -p 8500:8500 -p 53:53/udp -h server1 progrium/consul -server -bootstrap -ui-dir /ui -advertise $(dm ip nb-consul)

#docker-machine create -d virtualbox --swarm --swarm-master --swarm-discovery="consul://$(docker-machine ip nb-consul):8500" --engine-opt="cluster-store=consul://$(docker-machine ip nb-consul):8500" -engine-opt="cluster-advertise=eth1:2376" nb1
#docker-machine create -d virtualbox --swarm --swarm-discovery="consul://$(docker-machine ip nb-consul):8500" --engine-opt="cluster-store=consul://$(docker-machine ip nb-consul):8500" --engine-opt="cluster-advertise=eth1:2376" nb2
#docker-machine create -d virtualbox --swarm --swarm-discovery="consul://$(docker-machine ip nb-consul):8500" --engine-opt="cluster-store=consul://$(docker-machine ip nb-consul):8500" --engine-opt="cluster-advertise=eth1:2376" nb3

update_docker_host(){
# clear existing docker.local entry from /etc/hosts
sudo sed -i "/"${1}"\.local$/d" /etc/hosts
# get ip of running machine
export DOCKER_IP="$(docker-machine ip $1)" && sudo /bin/bash -c "echo \"${DOCKER_IP} $1.local\" >> /etc/hosts"
}

alias update-docker-host=update_docker_host
update-docker-host nb1
update-docker-host nb2
update-docker-host nb3
update-docker-host nb-consul