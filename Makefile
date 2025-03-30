# cli

.PHONY: run
run: bin/main
	cd bin && npx dotenv-cli -e ../.env -- ./main

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

	# Tailwind
	npm install

# targets

.PHONY: bin/main
bin/main: bin/tailwind.css
	go build -v -o bin/main .

.PHONY: bin/tailwind.css
bin/tailwind.css: input.css
	npx @tailwindcss/cli -i input.css -o bin/tailwind.css
