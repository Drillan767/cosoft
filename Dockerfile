FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o server ./slack/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
RUN mkdir -p /data
ENV DB_PATH=/data/database.db
VOLUME ["/data"]
EXPOSE 8080
CMD ["./server"]