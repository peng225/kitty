GOLANGCI_LINT_VERSION := v2.12.1

.PHONY: install
install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -count=1 -v ./...
