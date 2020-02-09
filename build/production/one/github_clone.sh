#!/bin/sh
echo ""
echo "  -----------------------------------"
echo "  -----$1 Firewall set-----"
echo "  -----------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x github_clone.sh && ./github_clone.sh 1

ip=$1   # the IP of this node

if [ -z "${ip}" ]; then
    chmod +x ./../print_error.sh && ./../print_error.sh \
        "$0 <ip of the node>"
    exit 1
fi

ssh root@$ip "
git clone https://github.com/go-park-mail-ru/2019_1_Escapade.git
cd 2019_1_Escapade
git checkout develop"