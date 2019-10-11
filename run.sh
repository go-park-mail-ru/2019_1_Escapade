#!/bin/sh
echo ""
echo "  -----------------------------"
echo "  ------docker compose up------"
echo "  -----------------------------"
echo ""
trap 'echo " stop" ' INT TERM

sudo docker-compose rm -v -f
sudo docker-compose up --scale ...