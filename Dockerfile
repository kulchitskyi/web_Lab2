FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .
RUN go build -o server ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/api ./api
COPY --from=builder /app/config.json ./config.json
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./server"]
