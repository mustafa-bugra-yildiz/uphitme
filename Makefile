.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: tools
tools:
	# DB migrations
	go install github.com/pressly/goose/v3/cmd/goose@latest
