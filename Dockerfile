FROM golang:1.15-alpine AS builder

WORKDIR /app/

# Get dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build the executable
COPY . .
RUN GOOS=linux go build -o main cmd/main.go

FROM alpine:latest as prod
WORKDIR /app/
COPY --from=builder /app/main .
CMD ["./main"]
