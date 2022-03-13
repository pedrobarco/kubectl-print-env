BINDIR      := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin
BINNAME     ?= kubectl-env

GOBIN         = $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN         = $(shell go env GOPATH)/bin
endif
ARCH          = $(shell uname -p)

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=

# Rebuild the binary if any of these files change
SRC := $(shell find . -type f -name '*.go' -print) go.mod go.sum

# Required for globs to work correctly
SHELL      = /usr/bin/env bash

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe | cut -c2-)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

VERSION_METADATA = unreleased

ifdef VERSION
	BINARY_VERSION = $(VERSION)
	VERSION_METADATA =
endif
BINARY_VERSION ?= $(GIT_TAG)

LDFLAGS += -X kubectl-env/internal/version.version=$(BINARY_VERSION)
LDFLAGS += -X kubectl-env/internal/version.metadata=$(VERSION_METADATA)
LDFLAGS += -X kubectl-env/internal/version.gitCommit=$(GIT_COMMIT)
LDFLAGS += $(EXT_LDFLAGS)

.PHONY: all
all: build

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)/$(BINNAME)' ./cmd/kubectl-env

# ------------------------------------------------------------------------------
#  install

.PHONY: install
install: build
	@install "$(BINDIR)/$(BINNAME)" "$(INSTALL_PATH)/$(BINNAME)"

# ------------------------------------------------------------------------------
#  test

.PHONY: test
test: build
ifeq ($(ARCH),s390x)
test: TESTFLAGS += -v
else
test: TESTFLAGS += -race -v
endif
test: test-unit

test-unit:
	@echo
	@echo "==> Running unit tests <=="
	GO111MODULE=on go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

# ------------------------------------------------------------------------------
# clean

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)'

# ------------------------------------------------------------------------------
# info

.PHONY: info
info:
	@echo "Version:           $(VERSION)"
	@echo "Git Tag:           $(GIT_TAG)"
	@echo "Git Commit:        $(GIT_COMMIT)"
	@echo "Git Tree State:    $(GIT_DIRTY)"
