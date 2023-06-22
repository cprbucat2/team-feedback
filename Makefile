.PHONY: build run-app tidy unittest vet fmt fmt-fix install-golangci lint \
	test clean vendor cover

export tf_VERSION := 0.0.0

export GO ?= go
export GOBUILDFLAGS += -trimpath -mod=readonly \
	-ldflags "-X github.com/cprbucat2/team-feedback.Version=$(tf_VERSION)"

GOTESTFLAGS ?= -race

SUBDIRS := app

.PHONY: $(SUBDIRS)
.PHONY: build-$(SUBDIRS) tidy-$(SUBDIRS) vet-$(SUBDIRS) lint-$(SUBDIRS) \
	unittest-$(SUBDIRS) clean-$(SUBDIRS) vendor-$(SUBDIRS)

build: build-$(SUBDIRS)

build-$(SUBDIRS):
	$(MAKE) -C $(@:build-%=%) build

run-app: build-app
	cd app && \
	./tf-server

tidy: tidy-$(SUBDIRS)
tidy-$(SUBDIRS):
	cd $(@:tidy-%=%) && go mod tidy

unittest: unittest-$(SUBDIRS)

unittest-$(SUBDIRS) $(SUBDIRS)/coverage.out:
	{ [ "$@" = "unittest-$(@:unittest-%=%)" ] && cd $(@:unittest-%=%) || \
	 cd $(@:%/coverage.out=%); } && \
	$(GO) test $(GOTESTFLAGS) -coverprofile=coverage.out .

vet: vet-$(SUBDIRS)
vet-$(SUBDIRS):
	cd $(@:vet-%=%) && go vet

fmt:
	gofmt -l $$(find . -type f -name '*.go' | grep -vF "/vendor/")
	@[ "`gofmt -l $$(find . -type f -name '*.go' | grep -v "/vendor/")`" = "" ]

fmt-fix:
	gofmt -w $$(find . -type f -name '*.go' | grep -vF "/vendor/")

install-golangci:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1

lint: lint-$(SUBDIRS)
lint-$(SUBDIRS):
	cd $(@:lint-%=%) && \
	$$($(GO) env GOPATH)/bin/golangci-lint run

test: vet fmt lint unittest

clean: clean-$(SUBDIRS)

clean-$(SUBDIRS):
	$(MAKE) -C $(@:clean-%=%) clean

vendor: vendor-$(SUBDIRS)
	cd $(@:vendor-%=%) && $(GO) mod vendor

cover: $(SUBDIRS)/coverage.out
	go tool cover -func=$<

coverage.html: coverage.out
	go tool cover -o $@ -html $<
