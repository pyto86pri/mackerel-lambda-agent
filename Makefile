BIN := mackerel-agent-lambda
VERSION := $(shell aws lambda list-layer-versions --layer-name mackerel-agent-lambda --query "LayerVersions[0].Version")
REVISION := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	sam build

BUILD_LDFLAGS := "\
	-X main.version=$(shell expr $(VERSION) + 1) \
	-X main.revision=$(REVISION)"

.PHONY: build-Layer
build-Layer:
	GOOS=linux go build -ldflags=$(BUILD_LDFLAGS) -o $(ARTIFACTS_DIR)/extensions/$(BIN) ./cmd/