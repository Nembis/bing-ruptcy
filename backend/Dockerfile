FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api/main.go

FROM golang:1.25.0 AS dev

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY . .

RUN go mod download

CMD ["air"]

FROM alpine:3.18.12 AS prod

RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 4000

CMD ["./main"]