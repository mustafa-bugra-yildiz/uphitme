.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: mocks
mocks:
	go generate ./...

.PHONY: tools
tools:
	# DB migrations
	go install github.com/pressly/goose/v3/cmd/goose@latest

	# Mock testing
	go install go.uber.org/mock/mockgen@latest
