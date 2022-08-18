# ########################################################## #
# Makefile for Golang Project
# Includes cross-compiling, installation, cleanup
# ########################################################## #

# Default Goal of the makefile is to show the help
.DEFAULT_GOAL := help

# Sets default shell to Bash
SHELL := /bin/bash

# Check for required command tools to build or stop immediately
#EXECUTABLES = git go awk docker golangci-lint staticcheck
#K := $(foreach exec,$(EXECUTABLES),\
 #       $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Setting up variables for the make targets
COMMIT_USER=$(shell git log -1 --pretty=format:'%an')
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD=`git rev-parse HEAD`
BINARY=podchaosmonkey
OUT_DIR=$(ROOT_DIR)/bin
# Setup linker flags option for build that interoperate with variable names in src code

# Indicates that the following targets have no physical files
.PHONY: all clean build lint staticcheck golangci-lint help test unit-test integration-test help docker-cleanup proto
clean: ## Clean the build artifacts
	rm -rf $(OUT_DIR)/$(BINARY)
build-static: ## Build the binary statically
	cd $(ROOT_DIR) && go mod tidy -v
	cd $(ROOT_DIR) && go build -ldflags "-w -s -X main.Version=$(VERSION) -X main.Build=$(BUILD)" -a -installsuffix cgo -o $(OUT_DIR)/$(BINARY)
build: ## Build the binary with dynamic linking
	cd $(ROOT_DIR) && go mod tidy -v
	cd $(ROOT_DIR) && go build -o $(OUT_DIR)/$(BINARY)
clean-build: clean build ## Clean the build artifacts & build dyamically
test: ## Run the unit tests
	go test -v -cover ./...
docker: ## build the docker image
	docker build -t $(BINARY):latest .
helm-test-workload: ## Install the test workload with helm
	helm upgrade --install test ./charts/testworkload --namespace workloads --create-namespace
helm-app: ## Install pod chaos monkey with helm
	helm upgrade --install test ./charts/podchaosmonkey --namespace workloads --create-namespace
helm: helm-test-workload helm-app ## Install the test-workload and pod chaos monkey with helm

all: clean test build ## Cleans runs, tests and builds the binary

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)