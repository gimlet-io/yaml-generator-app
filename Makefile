GOFILES = $(shell find . -type f -name '*.go' -not -path "./.git/*")
LDFLAGS = '-s -w -extldflags "-static" -X github.com/gimlet-io/gimlet-cli/pkg/version.Version='${VERSION}

.PHONY: format
.PHONY: build-yaml-generator-app dist-yaml-generator-app

format:
	@gofmt -w ${GOFILES}

build-yaml-generator-app:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o build/yaml-generator-app github.com/gimlet-io/yaml-generator-app/cmd
