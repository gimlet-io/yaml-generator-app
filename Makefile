GOFILES = $(shell find . -type f -name '*.go' -not -path "./.git/*")
LDFLAGS = '-s -w -extldflags "-static" -X github.com/gimlet-io/gimlet-cli/pkg/version.Version='${VERSION}

.PHONY: format
.PHONY: build-yaml-generator-app dist-yaml-generator-app

format:
	@gofmt -w ${GOFILES}

build-yaml-generator-app:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o build/yaml-generator-app github.com/gimlet-io/yaml-generator-app/cmd

dist-yaml-generator-app:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/yaml-generator-app-linux-x86_64 github.com/gimlet-io/yaml-generator-app/cmd
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/yaml-generator-app-darwin-x86_64 github.com/gimlet-io/yaml-generator-app/cmd
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/yaml-generator-app-darwin-arm64 github.com/gimlet-io/yaml-generator-app/cmd
