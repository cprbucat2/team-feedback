.PHONY: build clean

build: tf-server

tf-server: app.go
	$(GO) build $(GOBUILDFLAGS) -o $@ .

clean:
	rm -rf tf-server
