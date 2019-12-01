#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------- Tag images -------"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x run.sh && ./run.sh

docker push smartphonejava/api:latest
docker push smartphonejava/auth:latest
docker push smartphonejava/front:latest
docker push smartphonejava/chat:latest
docker push smartphonejava/game:latest
docker push smartphonejava/pg:latest
docker push smartphonejava/pg_ery:latest
docker push smartphonejava/consul:latest

















