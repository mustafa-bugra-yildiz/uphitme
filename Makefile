# cli

.PHONY: run
run:
	rm -rf bin
	make bin/main
	cd bin && npx dotenv-cli -e ../.env -- ./main

.PHONY: coverage
coverage: bin/main
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: tools
tools:
	# DB migrations
	go install github.com/pressly/goose/v3/cmd/goose@latest

	# Mock testing
	go install go.uber.org/mock/mockgen@latest

	# Tailwind
	npm install

# targets

bin/main: bin/tailwind.css
	go generate ./...
	go build -v -o bin/main .

bin/tailwind.css: input.css
	npx @tailwindcss/cli -i input.css -o bin/tailwind.css
