.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

APP=stevedore
SRC_PACKAGES=$(shell go list -mod=vendor ./... | grep -v "vendor" | grep -v "swagger")
VERSION?=2.3
HELM_VERSION?=2.16.9
BUILD?=$(shell git describe --always --dirty 2> /dev/null)
GOLINT:=$(shell command -v golint 2> /dev/null)
APP_EXECUTABLE="./out/$(APP)"
RICHGO=$(shell command -v richgo 2> /dev/null)
GOMETA_LINT=$(shell command -v golangci-lint 2> /dev/null)
GOLANGCI_LINT_VERSION=v1.27.0
GO111MODULE=on
SHELL=bash -o pipefail
DEFAULT_HELM_REPO_NAME?="chartmuseum"
BUILD_ARGS="-s -w -X main.version=$(VERSION) -X main.helmVersion=$(HELM_VERSION) -X main.build=$(BUILD) -X github.com/gojek/stevedore/cmd/repo.defaultHelmRepoName=$(DEFAULT_HELM_REPO_NAME)"

ifeq ($(GOMETA_LINT),)
	GOMETA_LINT=$(shell command -v $(PWD)/bin/golangci-lint 2> /dev/null)
endif

ifeq ($(RICHGO),)
	GO_BINARY=go
else
	GO_BINARY=richgo
endif

ifeq ($(BUILD),)
	BUILD=dev
endif

ifdef CI_COMMIT_SHORT_SHA
	BUILD=$(CI_COMMIT_SHORT_SHA)
endif

all: setup build

ci: setup-common build-common

ensure-build-dir:
	mkdir -p out

build-deps: ## Install dependencies
	go get
	go mod tidy
	go mod vendor

update-deps: ## Update dependencies
	go get -u

compile: compile-app ## Compile stevedore

compile-app: ensure-build-dir
	$(GO_BINARY) build -mod=vendor -ldflags $(BUILD_ARGS) -o $(APP_EXECUTABLE) ./main.go

install: ## Install stevedore
	go install -mod=vendor -ldflags $(BUILD_ARGS)

compile-linux: ensure-build-dir ## Compile stevedore for linux
	GOOS=linux GOARCH=amd64 $(GO_BINARY) build -mod=vendor -ldflags $(BUILD_ARGS) -o $(APP_EXECUTABLE) ./main.go

build: fmt build-common ## Build the application

build-common: vet lint-all test compile

compress: compile ## Compress the binary
	upx $(APP_EXECUTABLE)

fmt:
	GOFLAGS="-mod=vendor" $(GO_BINARY) fmt $(SRC_PACKAGES)

vet:
	$(GO_BINARY) vet -mod=vendor $(SRC_PACKAGES)

setup-common:
ifeq ($(GOMETA_LINT),)
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(GOLANGCI_LINT_VERSION)
endif
ifeq ($(GOLINT),)
	GO111MODULE=off $(GO_BINARY) get -u golang.org/x/lint/golint
endif

setup-richgo:
ifeq ($(RICHGO),)
	GO111MODULE=off $(GO_BINARY) get -u github.com/kyoh86/richgo
endif

setup: setup-richgo setup-common ensure-build-dir ## Setup environment

lint-all: lint setup-common
	$(GOMETA_LINT) run --modules-download-mode=vendor --timeout 2m0s

lint:
	./scripts/lint $(SRC_PACKAGES)

test-all: test test.integration

test: ensure-build-dir ## Run tests
	ENVIRONMENT=test $(GO_BINARY) test -mod=vendor $(SRC_PACKAGES) -race -coverprofile ./out/coverage -short -v | grep -viE "start|no test files"

test.integration: ensure-build-dir ## Run integration tests
	ENVIRONMENT=test $(GO_BINARY) test -mod=vendor $(SRC_PACKAGES) -tags=integration -short -v | grep -viE "start|no test files"

test-cover-html: ## Run tests with coverage
	mkdir -p ./out
	@echo "mode: count" > coverage-all.out
	$(foreach pkg, $(SRC_PACKAGES),\
	ENVIRONMENT=test $(GO_BINARY) test -mod=vendor -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	$(GO_BINARY) tool -mod=vendor cover -html=coverage-all.out -o out/coverage.html

dev-docker-build: ## Build stevedore server docker image with local chartmuseum repo added to it
	docker build --build-arg BUILD_MODE="dev" -f docker/stevedore/Dockerfile -t local-stevedore .
	@echo "To run:"
	@echo "$ docker run --rm -it -p 5443:5443 --name local-stevedore local-stevedore:latest"

dev-docker-run:  ## Run stevedore server in local docker container pointing to local(host-machine) chartmuseum
	docker run --rm -it -p 5443:5443 --name local-stevedore local-stevedore:latest

generate-doc: ## Generate swagger api doc
	go mod vendor
	GO111MODULE=off swagger generate spec -o swagger.json

generate-test-summary:
	ENVIRONMENT=test $(GO_BINARY) test -mod=vendor $(SRC_PACKAGES) -race -coverprofile ./out/coverage -short -v -json | grep -viE "start|no test files" | tee test-summary.json; \
    sed -i '' -E "s/^(.+\| {)/{/" test-summary.json; \
	passed=`cat test-summary.json | jq | rg '"Action": "pass"' | wc -l`; \
	skipped=`cat test-summary.json | jq | rg '"Action": "skip"' | wc -l`; \
	failed=`cat test-summary.json | jq | rg '"Action": "fail"' | wc -l`; \
	echo "Passed: $$passed | Failed: $$failed | Skipped: $$skipped"
