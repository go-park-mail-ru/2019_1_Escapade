#!/bin/sh
echo "  -----------------"
echo "  -[E]swagger docs-"
echo "  -----------------"
echo ""

# go get -u github.com/swaggo/swag/cmd/swag

echo "  1. Environment set"
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
echo "  2. Swagger documentation generate"

# https://github.com/swaggo/swag
swag init -g ../../../../cmd/services/api/main.go -d ../../../internal/services/api/handlers/ -o ../../../docs