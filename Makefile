GO=go
GOFLAGS=-race
DEV_BIN=bin
COV_PROFILE=cp.out

export CGO_ENABLED = 1

.DEFAULT_GOAL := build

fmt:
	$(GO) fmt ./...
.PHONY: fmt

vet: fmt
	$(GO) vet ./...
.PHONY: vet

lint: vet
	golint -set_exit_status=1 ./...
.PHONY: lint

test: lint
	$(GO) clean -testcache
	$(GO) test ./... -v
.PHONY: test

install: test
	$(GO) install ./...
.PHONY: install

build: test
	$(GO) build $(GOFLAGS) github.com/mdm-code/termcols/...
.PHONY: build

cover:
	$(GO) test -coverprofile=$(COV_PROFILE) -covermode=atomic ./...
	$(GO) tool cover -html=$(COV_PROFILE)
.PHONY: cover

clean:
	$(GO) clean github.com/mdm-code/termcols/...
	$(GO) mod tidy
	$(GO) clean -testcache
	rm -f $(COV_PROFILE)
.PHONY: clean
