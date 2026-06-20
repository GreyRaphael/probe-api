VERSION := 0.1.0
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build clean install

build:
	go build $(LDFLAGS) -o probe-api .

clean:
	rm -f probe-api probe-api.exe probe-api-*.tar.gz probe-api-*.zip checksums.txt

install:
	go install $(LDFLAGS) .

.PHONY: package
package: clean
	@mkdir -p dist
	@for target in "linux amd64" "linux arm64" "darwin amd64" "darwin arm64"; do \
		os=$$(echo $$target | cut -d' ' -f1); \
		arch=$$(echo $$target | cut -d' ' -f2); \
		echo "Building $$os/$$arch ..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o probe-api .; \
		tar czf "dist/probe-api-$${os}-$${arch}.tar.gz" probe-api; \
		rm -f probe-api; \
	done
	@echo "Building windows/amd64 ..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o probe-api.exe .
	@cd dist && zip probe-api-windows-amd64.zip ../probe-api.exe && rm -f ../probe-api.exe
	@cd dist && sha256sum probe-api-* > checksums.txt
	@echo "Done. Archives in dist/"
