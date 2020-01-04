#!/bin/sh
echo ""
echo "  -----------------------------------"
echo "  -----$1 Firewall set-----"
echo "  -----------------------------------"
echo ""
trap 'echo " stop" ' INT TERM

#$1 - the IP of this node

ssh root@$1 "
ufw allow 22/tcp
ufw allow 2376/tcp
ufw allow 2377/tcp
ufw allow 7946/tcp 
ufw allow 7946/udp 
ufw allow 4789/udp 
ufw allow 8786/tcp
ufw reload
yes | ufw enable
systemctl restart docker
ufw status verbose"