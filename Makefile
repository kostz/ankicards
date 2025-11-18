ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SHELL := /bin/sh

SOURCES = $(shell find $(ROOT_DIR) -name "*.go" -print )
SOURCE_DIRS = $(shell find $(ROOT_DIR) -d -print | grep -v . )
TESTS   = $(shell go list ./... | grep -v e2e )
COVERAGE_DIR ?= $(PWD)/coverage

export GO111MODULE = on

default: all

clean:
	rm -rf build

check: test lint checkfmt coverage

test:
	go test -race -v -failfast $(TESTS)

checkfmt:
	@[ -z $$(gofmt -l $(SOURCES)) ] || (echo "Sources not formatted correctly. Fix by running: make fmt" && false)

fmt: $(SOURCES)
	gofmt -s -w $(SOURCES)

lint:
	golint -set_exit_status $(SOURCE_DIRS)
	golangci-lint run

coverage:
	mkdir -p $(COVERAGE_DIR)
	go test -v $(TESTS) -coverpkg=./... -coverprofile=$(COVERAGE_DIR)/coverage.out
	go test -v $(TESTS) -coverpkg=./... -covermode=count -coverprofile=$(COVERAGE_DIR)/count.out fmt
	go tool cover -func=$(COVERAGE_DIR)/coverage.out
	go tool cover -func=$(COVERAGE_DIR)/count.out
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/index.html

build: $(SOURCES)
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" $(CMD) -o build/

extract-verbs:
	go run . extractVerbsFromImages

add-verb-examples:
	go run . addVerbExamples

ankicards:
	go run . makeAnkicards
