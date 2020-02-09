#!/bin/sh
echo ""
echo "  ------------------------------"
echo "  -----$1 Fix metrics------"
echo "  ------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x metrics.sh && ./metrics.sh 1

ip=$1   # the IP of this node

if [ -z "${ip}" ]; then
    chmod +x ./../print_error.sh && ./../print_error.sh \
        "$0 <ip of the node>"
    exit 1
fi

ssh root@$ip "
yes | sudo apt-get purge do-agent
curl -sSL https://repos.insights.digitalocean.com/install.sh | sudo bash
/opt/digitalocean/bin/do-agent --version
"