#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------- Tag images -------"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x run.sh && ./run.sh

docker tag 2019_1_escapade_api:latest localhost:5000/2019_1_escapade/api:latest
docker tag 2019_1_escapade_auth:latest localhost:5000/2019_1_escapade/auth:latest
docker tag 2019_1_escapade_front:latest localhost:5000/2019_1_escapade/front:latest
docker tag 2019_1_escapade_chat:latest localhost:5000/2019_1_escapade/chat:latest
docker tag 2019_1_escapade_game:latest localhost:5000/2019_1_escapade/game:latest
docker tag 2019_1_escapade_pg:latest localhost:5000/2019_1_escapade/pg:latest
docker tag 2019_1_escapade_pg_ery:latest localhost:5000/2019_1_escapade/pg_ery:latest
docker tag 2019_1_escapade_consul:latest localhost:5000/2019_1_escapade/consul:latest
docker tag 2019_1_escapade_consul-server:latest localhost:5000/2019_1_escapade/consul-server:latest
docker tag 2019_1_escapade_consul-agent:latest localhost:5000/2019_1_escapade/consul-agent:latest

docker push localhost:5000/2019_1_escapade/api:latest
docker push localhost:5000/2019_1_escapade/auth:latest
docker push localhost:5000/2019_1_escapade/front:latest
docker push localhost:5000/2019_1_escapade/chat:latest
docker push localhost:5000/2019_1_escapade/game:latest
docker push localhost:5000/2019_1_escapade/pg:latest
docker push localhost:5000/2019_1_escapade/pg_ery:latest
docker push localhost:5000/2019_1_escapade/consul:latest
docker push localhost:5000/2019_1_escapade/consul-server:latest
docker push localhost:5000/2019_1_escapade/consul-agent:latest

















