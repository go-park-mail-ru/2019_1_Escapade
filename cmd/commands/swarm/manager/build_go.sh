#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ---- Build golang bin ----"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x build_go.sh && ./build_go.sh

# ./cmd/services/auth/Dockerfile
#     context: ./../../../../
DOCKER_BUILDKIT=1 docker build --build-arg PROJECT=$PWD/../../../.. -t smartphonejava/auth -f ./../../../../cmd/services/auth/Dockerfile ./../../../../
# 0:43+