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

build: $(SUBDIRS:%=build-%)

$(SUBDIRS:%=build-%):
	$(MAKE) -C $(@:build-%=%) build

run-app: build-app
	cd app && \
	./tf-server

tidy: $(SUBDIRS:%=tidy-%)
$(SUBDIRS:%=tidy-%):
	cd $(@:tidy-%=%) && $(GO) mod tidy
	for package in $(SUBDIRS); do \
		$(GO) get github.com/cprbucat2/team-feedback/$$package; \
	done

unittest: $(SUBDIRS:%=unittest-%)

$(SUBDIRS:%=unittest-%) $(SUBDIRS:%=%/coverage.out):
	{ [ "$@" = "unittest-$(@:unittest-%=%)" ] && cd $(@:unittest-%=%) || \
	 cd $(@:%/coverage.out=%); } && \
	$(GO) test $(GOTESTFLAGS) -coverprofile=coverage.out ./...

vet: $(SUBDIRS:%=vet-%)
$(SUBDIRS:%=vet-%):
	cd $(@:vet-%=%) && $(GO) vet .

fmt:
	gofmt -l $$(find . -type f -name '*.go' | grep -vF "/vendor/")
	@[ "`gofmt -l $$(find . -type f -name '*.go' | grep -v "/vendor/")`" = "" ]

fmt-fix:
	gofmt -w $$(find . -type f -name '*.go' | grep -vF "/vendor/")

install-golangci:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1

lint: $(SUBDIRS:%=lint-%)
$(SUBDIRS:%=lint-%):
	cd $(@:lint-%=%) && \
	$$($(GO) env GOPATH)/bin/golangci-lint run

test: vet fmt lint unittest

clean: $(SUBDIRS:%=clean-%)

$(SUBDIRS:%=clean-%):
	$(MAKE) -C $(@:clean-%=%) clean

vendor: $(SUBDIRS:%=vendor-%)
$(SUBDIRS:%=vendor-%):
	cd $(@:vendor-%=%) && $(GO) mod vendor

cover: $(SUBDIRS:%=%/coverage.out)
	$(GO) tool cover $(SUBDIRS:%=-func=%/coverage.out)

coverage.html: $(SUBDIRS:%=%/coverage.out)
	$(GO) tool cover -o $@ $(SUBDIRS:%=-html=%/coverage.out)
