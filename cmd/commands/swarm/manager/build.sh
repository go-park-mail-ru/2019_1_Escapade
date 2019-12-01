#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------ Build images ------"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x build_go.sh && ./build_go.sh

# before -> 40+38+65+113
#

project=./../../../../
services=$project/cmd/services

DOCKER_BUILDKIT=1 docker build -t smartphonejava/auth -f $services/auth/Dockerfile $project 
DOCKER_BUILDKIT=1 docker build -t smartphonejava/api -f $services/api/Dockerfile $project
DOCKER_BUILDKIT=1 docker build -t smartphonejava/chat -f $services/chat/Dockerfile $project
DOCKER_BUILDKIT=1 docker build -t smartphonejava/game -f $services/game/Dockerfile $project
docker-compose build