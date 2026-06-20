VERSION := 0.1.0
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build clean install test

build:
	go build $(LDFLAGS) -o probe-api .

clean:
	rm -f probe-api probe-api.exe

install:
	go install $(LDFLAGS) .

# Cross-compile for common targets
.PHONY: build-all
build-all:
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/probe-api-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/probe-api-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/probe-api-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/probe-api-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/probe-api-windows-amd64.exe .

test:
	go test -v ./...
