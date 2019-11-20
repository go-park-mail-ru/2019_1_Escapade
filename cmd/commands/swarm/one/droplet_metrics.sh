#!/bin/sh
echo ""
echo "  ------------------------------"
echo "  -----Fix droplet metrics------"
echo "  ------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x droplet_metrics.sh && ./droplet_metrics.sh 1

#$1 - the IP of this node

ssh root@$1 "
yes | sudo apt-get purge do-agent
curl -sSL https://repos.insights.digitalocean.com/install.sh | sudo bash
/opt/digitalocean/bin/do-agent --version
"