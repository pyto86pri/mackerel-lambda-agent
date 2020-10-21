BIN := mackerel-lambda-agent
VERSION = $(shell aws lambda list-layer-versions --layer-name mackerel-lambda-agent --query "LayerVersions[0].Version")
REVISION = $(shell git rev-parse --short HEAD)

.PHONY: fmt
fmt:
	test -z $(shell gofmt -l **/*.go)

.PHONY: fmt-fix
fmt-fix:
	gofmt -l -w **/*.go

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	GOOS=linux go build -o $(BIN) ./cmd/

# following used with AWS SAM CLI
.PHONY: sam-build
sam-build:
	sam build

BUILD_LDFLAGS = "\
	-X main.version=$(shell expr $(VERSION) + 1) \
	-X main.revision=$(REVISION)"

.PHONY: build-Layer
build-Layer:
	GOOS=linux go build -ldflags=$(BUILD_LDFLAGS) -o $(ARTIFACTS_DIR)/extensions/$(BIN) ./cmd/