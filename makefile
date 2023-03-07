SHELL := /bin/bash

.PHONY: all build test deps deps-cleancache
 
 ## varable for `go`
GOCMD=go 
 ## variable  `build`
BUILD_DIR=build
# variable as build/bin
BINARY_DIR=$(BUILD_DIR)/bin  
# variable code-coverage
CODE_COVERAGE=code-coverage 

all: test build
	
 ## to create binary directory if its not available(-p)
${BINARY_DIR}:
	mkdir -p $(BINARY_DIR)

 ## first call the binary_dir ## next build go file from ./cmd/api all files to binary_dir
build : ${BINARY_DIR} 
	$(GOCMD) build -o ${BINARY_DIR} -v ./cmd/api 

 # to start the application
run:
	@echo "Welcome to my ecommerce"
	$(GOCMD) run ./cmd/api

 # to test all tests in current and sub modlues
test:
	$(GOCMD) test ./... -cover

 # to test the tests and store on variable code_coverage and show as an html
test-coverage:
	$(GOCMD) test ./... -coverprofiile=$(CODE_COVERAGE).out
	$(GOCMD) tool cover -html=$(CODE_COVERAGE).out

# to install dependencies packges latest version if its not in local package
deps: 
	$(GOCMD) get -u -t -d -v ./...
	#remove un used dependencies
	$(GOCMD) mod tidy # 
	# create a vendor file in local 
	$(GOCMD) mod vendor

 # to clean cache in the module
dps-cleancache:
	$(GOCMD) clean -modcache

 ## Generate wire_gen.go
wire:
	cd pkg/di && wire

## Generate swagger docs
swag: 
	swag init -g pkg/http/handler/user.go -o ./cmd/api/docs
 
## Display this help screen
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'