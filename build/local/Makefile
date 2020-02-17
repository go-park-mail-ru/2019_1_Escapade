.PHONY: init run stop status build restart
.PHONY: restart-with-rebuild rebuild-network clear
,PHONY: drain active help

STACK=local
NETWORK=traefik
VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

help: ## - Show help message
	@printf "\n usage: make [target]\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:	## - Build services
	@printf "\n\033[32m Build images and run \n\033[0m"
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build
	@docker-compose push

run:	## - Run services
	@printf "\n\033[32m Build images and run \n\033[0m"
	@docker stack rm ${STACK}
	@env $(cat .env | grep ^[A-Z] | xargs) docker stack deploy -c docker-compose.yaml ${STACK}

stop:	## - Stop services
	@printf "\n\033[32m Stop services \n\033[0m"
	@docker stack rm ${STACK}

env:	## - View flags and environment variables accepted by each service
	@printf "\n\033[32m Api service: flags and environment variables  \n\033[0m"
	@go run ./../../cmd/api/main.go -h

clean:	## - Clean dangling containers 
	@printf "\n\033[32m Clean dangling containers \n\033[0m"
	@./../scripts/machine/clear.sh

run-with-rebuild:	## - Run services
	@printf "\n\033[32m Build images and run \n\033[0m"
	@docker stack rm ${STACK}
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build
	@docker-compose push
	@docker stack deploy -c docker-compose.yaml ${STACK}

rebuild-network:	## - Remove and add services network
	@printf "\n\033[32m Remove and add services network \n\033[0m"
	@docker network rm ${NETWORK}
	@docker network create -d overlay --subnet 10.10.8.0/24 --attachable ${NETWORK}
