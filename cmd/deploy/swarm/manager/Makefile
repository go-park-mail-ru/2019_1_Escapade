.PHONY: help init run stop status build restart
.PHONY: restart-with-rebuild rebuild-network clear
,PHONY: logs drain active

VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

help: ## - Show help message
	@printf "\n usage: make [target]\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init:	## - Install libs and bins, create swap
	@printf "\n\033[32m Install libs and bins, create swap\n\033[0m"
	@chmod +x install.sh && chmod +x swap.sh
	@./swap.sh && ./install.sh

run:	## - Run services in swarm 
	@printf "\n\033[32m Run services in swarm \n\033[0m"
	@docker stack deploy -c docker-swarm.yaml app

stop:	## - Stop services  
	@printf "\n\033[32m Stop services\n\033[0m"
	@docker stack rm app

status:	## - Status of services  
	@printf "\n\033[32m Status of services \n\033[0m"
	@docker service ls

build:	## - Build services
	@printf "\n\033[32m Build images and push them to dockerhub \n\033[0m"
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build
	@docker-compose push

restart:	## - Restart services
	@printf "\n\033[32m Restart services in swarm \n\033[0m"
	@docker stack rm app
	@docker stack deploy -c docker-swarm.yaml app

restart-with-rebuild:	## - Stop, rebuild and start services
	@printf "\n\033[32m Stop, rebuild and start services \n\033[0m"
	@docker stack rm app
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build 
	# --parallel
	@docker-compose push
	@docker stack deploy -c docker-swarm.yaml app

rebuild-network:	## - Remove and add services network
	@printf "\n\033[32m Remove and add services network \n\033[0m"
	@docker network rm traefik
	@docker network create -d overlay --subnet 10.10.8.0/24 --attachable traefik

clean:	## - Clean dangling containers 
	@printf "\n\033[32m Clean dangling containers \n\033[0m"
	@./clear.sh

active:	## - Mark node as active. Call me as 'make active node=<name>'
	@./node_active.sh ${node}
	@docker node ls

drain:	## - Mark node as drain. Call me as 'make drain node=<name>'
	@./node_drain.sh ${node}
	@docker node ls
