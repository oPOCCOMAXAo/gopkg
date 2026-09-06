GOLANGCI_LINT := ./bin/golangci-lint

fmt:
	$(GOLANGCI_LINT) fmt

lint:
	$(GOLANGCI_LINT) run

install-dev-tools:
	@mkdir -p ./bin
	@curl -sSfL https://golangci-lint.run/install.sh | sh -s v2.13.1

test:
	go test -v ./...
