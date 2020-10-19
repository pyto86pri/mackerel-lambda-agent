.PHONY: build
build:
	sam build

.PHONY: build-Layer
build-Layer:
	GOOS=linux go build -o $(ARTIFACTS_DIR)/extensions/mackerel-agent-lambda ./cmd/