.PHONY: help build run stop

export VERSION := $(if $(TAG),$(TAG),$(if $(BRANCH_NAME),$(BRANCH_NAME),$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)))

help: ## List all available targets with help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build:
	@cd ./rmq && make build
	@cd ./kafka && make build

dev-infra-up:
	@docker-compose up -d rabbitmq zookeeper kafka
	@sleep 20

run: dev-infra-up ## Run service in develop environment
	@docker-compose up app-rmq app-kafka

stop: ## Stop develop environment
	@docker-compose down
