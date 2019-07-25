#!/bin/sh
echo "  -----------------"
echo "  -[E]swagger docs-"
echo "  -----------------"
echo ""

echo "  1. Environment set"
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
echo "  2. Swagger documentation generate"
swag init -g main.go -d ../api -o ../docs