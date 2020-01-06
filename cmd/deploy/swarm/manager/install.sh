#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ---- Install packages ----"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x install.sh && ./install.sh

echo "1. Install docker"
place=$(whereis docker | grep bin)
if [ "$?" = "0" ]; then
    echo "✔ " $place
else
    yes | apt-get update
    yes | apt install apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    yes | sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
    yes | sudo apt update
    apt-cache policy docker-ce
    yes | sudo apt install docker-ce
    sudo systemctl status docker
fi

echo "\n2. Install docker-compose"
place=$(whereis docker-compose | grep bin)
if [ "$?" = "0" ]; then
    echo "✔ " $place
else
    sudo apt-get -y install python-pip
    sudo pip install docker-compose
fi

echo "\n3. Install dive for docker images check"
place=$(whereis dive | grep bin)
if [ "$?" = "0" ]; then
    echo "✔ " $place
else
    wget https://github.com/wagoodman/dive/releases/download/v0.9.1/dive_0.9.1_linux_amd64.deb
    sudo apt install ./dive_0.9.1_linux_amd64.deb
fi

echo "\n4. Install golang"
export PATH=$PATH:/usr/local/go/bin
place=$(whereis go | grep bin)
if [ "$?" = "0" ]; then
    echo "✔ " $place
else
    VERSION=1.13.4
    OS=linux
    ARCH=amd64
    sudo curl -O https://storage.googleapis.com/golang/go$VERSION.$OS-$ARCH.tar.gz
    tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
    rm go$VERSION.$OS-$ARCH.tar.gz
    export PATH=$PATH:/usr/local/go/bin
fi