#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ---Consul server launch---"
echo "  --------------------------"
echo ""
#chmod +x consul.sh && ./consul.sh

consul agent -server -ui -bootstrap -data-dir /tmp/consul

