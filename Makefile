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
	GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/linux/amd64/yaml-generator-app github.com/gimlet-io/yaml-generator-app/cmd
	GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/linux/arm64/yaml-generator-app github.com/gimlet-io/yaml-generator-app/cmd