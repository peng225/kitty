.PHONY: lint
lint:
	go tool golangci-lint run

.PHONY: test
test:
	go test -count=1 -v ./...
