ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

	COPY go.mod go.sum ./
	RUN go mod download && go mod verify

	COPY . .
	RUN go build -v -o uphitme .

FROM scratch
WORKDIR /app
COPY --from=builder /src/uphitme /app/uphitme
CMD ["/app/uphitme"]