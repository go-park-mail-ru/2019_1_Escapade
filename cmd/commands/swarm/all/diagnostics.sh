#!/bin/sh
echo ""
echo "  ------------------------------"
echo "  ----Nodes check ports open----"
echo "  ------------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x prepare.sh && ./prepare.sh

diagnostics() {
    script="./one/diagnostics.sh"
    chmod +x $script && \
    $script $1 $2
}

# $1 - $5 - адреса серверов
diagnostics $2 $1 && \
diagnostics $1 $2 && \
diagnostics $1 $3 && \
diagnostics $1 $4 && \
diagnostics $1 $5 

# script1="$script {$2} {$1}"
# script2="$script $1 $2"
# script3="$script $1 $3"
# script4="$script $1 $4"
# script5="$script $1 $5"

# $do $script1 $script2 $script3 $script4 $script5


# exec="/one/diagnostics.sh"
# chmod +x .$exec && \


# echo "  1. Check open ports on machines"

# echo "  1.1."$addr1 && \
#  .$exec $2 $1 && \

# echo "  1.2."$addr2 && \
#  ./$exec $1 $2 && \

# echo "  1.3."$addr3 && \
#  ./$exec $1 $3 && \

# echo "  1.4."$addr4 && \
#  ./$exec $1 $4 && \

# echo "  1.5."$addr5 && \
#  ./$exec $1 $5
