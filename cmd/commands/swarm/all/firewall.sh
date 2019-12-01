#!/bin/sh
echo ""
echo "  ----------------------------"
echo "  -----Nodes Firewall set-----"
echo "  ----------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x prepare.sh && ./prepare.sh

firewall() {
    script="./one/firewall.sh"
    chmod +x $script && \
    $script $1 $2
}

# $1 - $5 - адреса серверов
firewall $1 && \
firewall $2 && \
firewall $3 && \
firewall $4 && \
firewall $5 


# echo "  1. Set firewall"

#  echo "  1.1."$addr1 && \
#  .$exec $1 && \

# echo "  1.2."$addr2 && \
#  ./$exec $2 && \

# echo "  1.3."$addr3 && \
#  ./$exec $3 && \

# echo "  1.4."$addr4 && \
#  ./$exec $4 && \

# echo "  1.5."$addr5 && \
#  ./$exec $5
 
