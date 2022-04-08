SHELL := bash# we want bash behaviour in all shell invocations

GIT_REVISION := $(shell git rev-parse --short HEAD)
ERR = echo ${TIME} ${RED}[FAIL]${CNone}

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences
BOLD := \033[1m
NORMAL := \033[0m
GREEN := \033[1;32m

OK = echo [ OK ]

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
	CGO_ENABLED=0 go build -o bin/hln -ldflags '-s -w -X github.com/h8r-dev/heighliner/pkg/version.Revision=$(GIT_REVISION)' ./cmd/main.go
	@echo "Saved to bin/hln"


.PHONY: test
test: check-diff unit-test-core # Run tests
	@$(OK) unit-tests pass

reviewable: fmt vet lint staticcheck # Make your PR ready to review
	go mod tidy

check-diff: reviewable # Execute auto-gen code commands and ensure branch is clean
	git --no-pager diff
	git diff --quiet || ($(ERR) please run 'make reviewable' to include all changes && false)
	@$(OK) branch is clean

vet:
	go vet ./...

unit-test-core:
	go test ./...

lint:
	GOBIN=$(GOBIN) ./scripts/ci/install_golangci.sh
	GOLANGCILINT=$(shell which golangci-lint)
	$(GOLANGCILINT) run ./...

staticcheck:
	./scripts/ci/install_staticcheck.sh
	STATICCHECK=$(shell which staticcheck)
	$(STATICCHECK) ./...

# Run go fmt against code
fmt:
	go fmt ./...

