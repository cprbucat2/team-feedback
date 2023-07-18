.PHONY: build run tidy unittest vet fmt fmt-fix install-golangci lint test clean vendor cover

VERSION := 0.0.0

GOFLAGS = -ldflags "-X github.com/cprbucat2/team-feedback.Version=$(VERSION)" -trimpath -mod=readonly
GOTESTFLAGS ?= -race

build: tf-server

run: build
	./tf-server

tidy:
	go mod tidy

tf-server: app.go go.sum go.mod
	go build $(GOFLAGS) -o $@ .

unittest coverage.out:
	go test $(GOTESTFLAGS) -coverprofile=coverage.out ./...

vet:
	go vet ./...

fmt:
	gofmt -l $$(find . -type f -name '*.go' | grep -v "/vendor/")
	@[ "`gofmt -l $$(find . -type f -name '*.go' | grep -v "/vendor/")`" = "" ]

fmt-fix:
	gofmt -w $$(find . -type f -name '*.go' | grep -v "/vendor/")

install-golangci:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1 run

test: vet fmt lint unittest

clean:
	rm -rf tf-server
	rm -rf coverage.out coverage.html

vendor:
	go mod vendor

cover: coverage.out
	go tool cover -func $<

coverage.html: coverage.out
	go tool cover -o $@ -html $<
