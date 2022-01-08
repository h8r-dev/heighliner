SHELL := bash# we want bash behaviour in all shell invocations

# https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences
BOLD := \033[1m
NORMAL := \033[0m
GREEN := \033[1;32m

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
hln: # build a client binary
	CGO_ENABLED=0 go build -o bin/hln -ldflags '-s -w -X github.com/h8r-dev/heighliner/pkg/version.Revision=$(GIT_REVISION)' ./cmd/client/main.go

GIT_REVISION := $(shell git rev-parse --short HEAD)
.PHONY: server
server: # build a server binary
	CGO_ENABLED=0 go build -o bin/server '-s -w -X github.com/h8r-dev/heighliner/pkg/version.Revision=$(GIT_REVISION)' ./cmd/server/main.go -ldflags
