SHELL := bash
.ONESHELL:
MAKEFLAGS += --no-builtin-rules

.PHONY: run stop build lint help

export REPO := rebrain
export APP_NAME := kafka
export VERSION := $(if $(TAG),$(TAG),$(if $(BRANCH_NAME),$(BRANCH_NAME),$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)))

NOCACHE := $(if $(NOCACHE),"--no-cache")

help: ## List all available targets with help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test-short: ## Run unit tests
	@go test ./... -cover -short

lint: ## Run golangci-lint
	@golangci-lint run

build: ## Build develop containers
	@docker build ${NOCACHE} --pull -f ./build/Dockerfile -t ${REPO}/${APP_NAME}:${VERSION} .

run: dev-infra-up ## Run service in develop environment
	@docker-compose up app

test-long: ## Run unit tests long
	@make dev-infra-up && make dev-migration-up && make test-long-up; make stop

test-long-up:
	@docker-compose run --rm helper /bin/bash ./scripts/long_test_runner.sh
