.PHONY: build run test clean tidy unittest

VERSION := 0.0.0

GOFLAGS = -ldflags "-X github.com/cprbucat2/team-feedback.Version=$(VERSION)" -trimpath -mod=readonly
GOTESTFLAGS ?= -race

build: tf-server

run:
	go run .

tidy:
	go mod tidy

tf-server:
	go build $(GOFLAGS) -o $@ .

unittest:
	go test $(GOTESTFLAGS) -coverprofile=coverage.out ./...

vet:
	go vet ./...

lint:
	gofmt -l $$(find . -type f -name '*.go' | grep -v "/vendor/")
	[ "`gofmt -l $$(find . -type f -name '*.go' | grep -v "/vendor/")`" = "" ]

test: lint vet unittest

clean:
	rm -rf tf-server

vendor:
	go mod vendor
