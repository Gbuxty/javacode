FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /javacode ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /javacode .
COPY --from=builder /app/.env ./
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./javacode"]