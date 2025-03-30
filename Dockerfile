ARG GO_VERSION=1

# build tailwind
FROM node:22.14.0 AS tailwind-builder
WORKDIR /src
	COPY package.json package-lock.json ./
	RUN npm install

	COPY . .
	RUN npx @tailwindcss/cli -i input.css -o bin/tailwind.css

# test & build binary
FROM golang:${GO_VERSION}-alpine AS main-builder
WORKDIR /src
	RUN apk add make

	COPY go.mod go.sum ./
	RUN go mod download && go mod verify

	COPY --from=tailwind-builder /src/bin /src/bin

	COPY . .
	RUN go test -v ./...
	RUN go build -v -o uphitme .

# runtime image
FROM scratch
WORKDIR /app
COPY --from=main-builder /src/uphitme /app/uphitme
CMD ["/app/uphitme"]
