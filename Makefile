GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
MAIN_PATH=./cmd/tiny-socks
BINARY_NAME=tiny-socks
IMAGE_NAME=stdtom/tiny-socks
VERSION?=0.0.1
SERVICE_PORT?=6000
DOCKER_REGISTRY?=ghcr.io/
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor

all: help

## Build:
build: ## Build your project and put the output binary in out/bin/
	mkdir -p out/bin
	GO111MODULE=on $(GOCMD) build -mod vendor -o out/bin/$(BINARY_NAME) $(MAIN_PATH)

clean: ## Remove build related file
	rm -fr ./bin
	rm -fr ./out
	rm -f ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov yamllint-checkstyle.xml

vendor: ## Copy of all packages needed to support builds and tests in the vendor directory
	$(GOCMD) mod tidy
	$(GOCMD) mod vendor

## Format / Lint / Code Analysis
fmt: ## Fix source code with gofmt
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

lint: ## Check source code against linters
	golangci-lint run ./...

vet: ## Code analysis with go vet
	go vet ./...

## Test:
test: ## Run the tests of the project
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/AlekSi/gocov-xml
	GO111MODULE=off go get -u github.com/axw/gocov/gocov
	gocov convert profile.cov | gocov-xml > coverage.xml
endif

## Docker:
docker-build: ## Use the dockerfile to build the container
	docker build --rm --tag $(IMAGE_NAME) .

docker-release: ## Release the container with tag latest and version
	docker tag $(IMAGE_NAME) $(DOCKER_REGISTRY)$(IMAGE_NAME):latest
	docker tag $(IMAGE_NAME) $(DOCKER_REGISTRY)$(IMAGE_NAME):$(VERSION)
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(IMAGE_NAME):latest
	docker push $(DOCKER_REGISTRY)$(IMAGE_NAME):$(VERSION)

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
