SHELL := bash# we want bash behaviour in all shell invocations

GOLANGCILINT_VERSION=latest

# https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences
BOLD := \033[1m
NORMAL := \033[0m
GREEN := \033[1;32m

OK		= echo [ OK ]

XDG_CONFIG_HOME ?= $(CURDIR)/.config
export XDG_CONFIG_HOME
.DEFAULT_GOAL := help
HELP_TARGET_DEPTH ?= \#
.PHONY: help
help: # Show how to get started & what targets are available
	@printf "This is a list of all the make targets that you can run, e.g. $(BOLD)make hln$(NORMAL) - or $(BOLD)m hln$(NORMAL)\n\n"
	@awk -F':+ |$(HELP_TARGET_DEPTH)' '/^[0-9a-zA-Z._%-]+:+.+$(HELP_TARGET_DEPTH).+$$/ { printf "$(GREEN)%-20s\033[0m %s\n", $$1, $$3 }' $(MAKEFILE_LIST) | sort
	@echo

.PHONY: hln
hln: # build client binary
	CGO_ENABLED=0 go build -o bin/hln -ldflags '-s -w -X github.com/h8r-dev/heighliner/pkg/version.Revision=$(GIT_REVISION)' ./cmd/client/main.go
	@echo "Saved to bin/hln"

GIT_REVISION := $(shell git rev-parse --short HEAD)
.PHONY: server
server: # build server binary
	CGO_ENABLED=0 go build -o bin/server '-s -w -X github.com/h8r-dev/heighliner/pkg/version.Revision=$(GIT_REVISION)' ./cmd/server/main.go -ldflags

# Run tests
.PHONY: test
test: vet lint unit-test-core
	@$(OK) unit-tests pass

# Run go vet against code
vet:
	go vet ./...

lint: golangci
	$(GOLANGCILINT) run ./...

unit-test-core:
	go test ./...

.PHONY: golangci
golangci: gobin
ifneq ($(shell which golangci-lint),)
	@$(OK) golangci-lint is already installed
GOLANGCILINT=$(shell which golangci-lint)
else ifeq (, $(shell which $(GOBIN)/golangci-lint))
	@{ \
	set -e ;\
	echo 'installing golangci-lint-$(GOLANGCILINT_VERSION)' ;\
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCILINT_VERSION) ;\
	echo 'Successfully installed' ;\
	}
GOLANGCILINT=$(GOBIN)/golangci-lint
else
	@$(OK) golangci-lint is already installed
GOLANGCILINT=$(GOBIN)/golangci-lint
endif

.PHONY: gobin
gobin:
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif
