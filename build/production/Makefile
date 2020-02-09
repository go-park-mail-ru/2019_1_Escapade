.PHONY: help init run
VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

help: ## - Show help message
	@printf "\n usage: make [target]\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init:	## - Install libs and prepare servers from info.txt
	@printf "\n\033[32m Install docker, add servers to docker-machines and swarm, set firewall and verify them\n\033[0m"
	@chmod +x install.sh && chmod +x init.sh
	@./install.sh && ./init.sh

run:	## - To manage the process of building, starting and stopping services, go to manager node and open manager/Makefile
	@run "\n\033[32m To manage the process of building, starting and stopping services, go to manager node and open manager/Makefile\n\033[0m"
	@docker-machine ssh api1 "cd 2019_1_Escapade/cmd/deploy/swarm/manager && make help"
