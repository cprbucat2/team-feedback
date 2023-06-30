.PHONY: build clean

build: combine-html

combine-html: main.go rewrite.go
	$(GO) build $(GOBUILDFLAGS) -o $@ .

clean:
	rm -rf combine-html
