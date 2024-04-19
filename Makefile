PROJET_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

PROJET_NAME:=$(shell basename ${PROJET_PATH})
COMPOSE_FILE:=./scripts/docker-compose.yaml

.PHONY: build
build:
	@docker compose -p ${PROJET_NAME} -f ${COMPOSE_FILE} build app

.PHONY: up
up:
	@docker compose -p ${PROJET_NAME} -f ${COMPOSE_FILE} up -d

.PHONY: stop
stop:
	@docker compose -p ${PROJET_NAME} -f ${COMPOSE_FILE} stop

.PHONY: down
down:
	@docker compose -p ${PROJET_NAME} -f ${COMPOSE_FILE} down -v

.PHONY: log
log:
	@docker compose -p ${PROJET_NAME} -f ${COMPOSE_FILE} logs -f