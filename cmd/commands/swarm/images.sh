#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------- Tag images -------"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x run.sh && ./run.sh

docker tag 2019_1_escapade_api:latest smartphonejava/api:latest
docker tag 2019_1_escapade_auth:latest smartphonejava/auth:latest
docker tag 2019_1_escapade_front:latest smartphonejava/front:latest
docker tag 2019_1_escapade_chat:latest smartphonejava/chat:latest
docker tag 2019_1_escapade_game:latest smartphonejava/game:latest
docker tag 2019_1_escapade_pg:latest smartphonejava/pg:latest
docker tag 2019_1_escapade_pg_ery:latest smartphonejava/pg_ery:latest
docker tag 2019_1_escapade_consul:latest smartphonejava/consul:latest
docker tag 2019_1_escapade_consul-server:latest smartphonejava/consul-server:latest
docker tag 2019_1_escapade_consul-agent:latest smartphonejava/consul-agent:latest

docker push smartphonejava/api:latest
docker push smartphonejava/auth:latest
docker push smartphonejava/front:latest
docker push smartphonejava/chat:latest
docker push smartphonejava/game:latest
docker push smartphonejava/pg:latest
docker push smartphonejava/pg_ery:latest
docker push smartphonejava/consul:latest
docker push smartphonejava/consul-server:latest
docker push smartphonejava/consul-agent:latest

















