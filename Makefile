PLUGIN_ID      := com.joelmgallant.memes
PLUGIN_VERSION := 1.0.0
BUNDLE_NAME    := $(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz

.PHONY: all test dist clean

all: test dist

test:
	@test -z "$$(gofmt -l .)" || { echo "unformatted files:"; gofmt -l .; exit 1; }
	go test ./server/...

dist: clean
	mkdir -p dist/$(PLUGIN_ID)/server/dist
	cp plugin.json dist/$(PLUGIN_ID)/
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
		-ldflags="-s -w" \
		-o dist/$(PLUGIN_ID)/server/dist/plugin-linux-amd64 ./server
	cd dist && tar -czf $(BUNDLE_NAME) $(PLUGIN_ID)
	@echo "built dist/$(BUNDLE_NAME)"

clean:
	rm -rf dist
