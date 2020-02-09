#!/bin/bash
#chmod +x logs.sh && ./logs.sh

stack=$1
if [ -z "${stack}" ]; then
    echo 'Usage:'
    echo "$0 <stack>"
    exit 1
fi

place=$(whereis xpanes | grep bin)
if [ "$?" = "0" ]; then
    echo "✔ xpanes is in " $place
else
    yes | apt install software-properties-common
    yes | add-apt-repository ppa:greymd/tmux-xpanes
    yes | apt update
    yes | apt install tmux-xpanes
    echo "✔ xpanes is installed "
fi

# tmux -S /root/.cache/xpanes/socket.14153 attach-session -t xpanes-14153
some=$(docker stack services $stack --format '{{.Name}}' | \
    xpanes -c "docker service logs -f {}")

