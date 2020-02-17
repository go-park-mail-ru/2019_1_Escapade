#!/bin/sh
echo "  --------------------------------"
echo "  -----Documentation generate-----"
echo "  --------------------------------"
echo ""
#chmod +x proto.sh && ./proto.sh

export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

SWAG_VERSION="1.6.4"
SWAG_ARCH="1.6.4_Linux_x86_64"

echo "1. Install swag\n"
mkdir swag
wget https://github.com/swaggo/swag/releases/download/v$SWAG_VERSION/swag_$SWAG_ARCH.tar.gz
tar -C swag -xvf  swag_$SWAG_ARCH.tar.gz
echo "\n✔ installed \n"

# https://github.com/swaggo/swag
echo "2. Swagger documentation generate"
swag/swag init -g  ../cmd/services/api/main.go -d ../../internal/ -o ../../docs

echo "3. Delete swag"
rm -R swag
rm swag_$SWAG_ARCH.tar.gz