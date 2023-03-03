SHELL := /bin/bash

.PHONY: all build test deps deps-cleancache

GOCMD=go  ## varable for `go`
BUILD_DIR=build ## variable  `build`
BINARY_DIR=$(BUILD_DIR)/bin  # variable as build/bin
CODE_COVERAGE=code-coverage # variable code-coverage

all: # to test and build the application
	test build

${BINARY_DIR}: ## to create binary directory if its not available(-p)
	mkdir -p $(BINARY_DIR)

build : ## first call the binary_dir ## next build go file from ./cmd/api all files to binary_dir
	${BINARY_DIR} 
	$(GOCMD) build -o ${BINARY_DIR} -v ./cmd/api 

run: # to start the application
	$(GOCMD) run ./cmd/api

test: # to test all tests in current and sub modlues
	$(GOCMD) test ./... -cover

test-coverage: # to test the tests and store on variable code_coverage and show as an html
	$(GOCMD) test ./... -coverprofiile=$(CODE_COVERAGE).out
	$(GOCMD) tool cover -html=$(CODE_COVERAGE).out

deps: # to install dependencies packges latest version if its not in local package
	$(GOCMD) get -u -t -d -v ./...
	#remove un used dependencies
	$(GOCMD) mod tidy # 
	# create a vendor file in local 
	$(GOCMD) mod vendor


dps-cleancache: # to clean cache in the module
	$(GOCMD) clean -modcache

wire: ## Generate wire_gen.go
	cd pkg/di && wire

swag: ## Generate swagger docs
	swag init -g pkg/http/handler/user.go -o ./cmd/api/docs

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'