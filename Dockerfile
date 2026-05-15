FROM golang:1.23-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o uptime-worker .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/uptime-worker ./uptime-worker
COPY db ./db
EXPOSE 8081
CMD ["./uptime-worker"]
