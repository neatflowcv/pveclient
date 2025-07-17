.PHONY: build
build:
	go build -o pveclient cmd/pveclient/main.go

.PHONY: lint
lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2
	golangci-lint run

.PHONY: test
test:
	go test -v ./...