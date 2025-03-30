ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine AS main-builder
WORKDIR /src

	COPY go.mod go.sum ./
	RUN go mod download && go mod verify

	COPY . .
	RUN go build -v -o uphitme .

FROM node:22.14.0 AS tailwind-builder
WORKDIR /src

	COPY package.json package-lock.json ./
	RUN npm install

	COPY . .
	RUN npx @tailwindcss/cli -i input.css -o tailwind.css


FROM scratch
WORKDIR /app
COPY --from=main-builder     /src/uphitme      /app/uphitme
COPY --from=tailwind-builder /src/tailwind.css /app/tailwind.css
CMD ["/app/uphitme"]
