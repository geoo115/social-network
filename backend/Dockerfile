# First stage: build the Go binary
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install dependencies required by go-sqlite3
RUN apk --no-cache add build-base sqlite-dev

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all Go modules
RUN go mod download

# Copy the entire Go project to the working directory
COPY . .

# Set CGO_ENABLED to 1 to enable CGO
ENV CGO_ENABLED=1

# Set the working directory to the directory where main.go is located
WORKDIR /app/cmd

# Build the Go app
RUN go build -o main .

# Second stage: create a lightweight container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/cmd/main .

# Copy the .env file (if needed) from the builder stage
COPY --from=builder /app/cmd/.env .env

# Install SQLite library for runtime
RUN apk --no-cache add sqlite-libs

# Expose the application port
EXPOSE 8000

# Command to run the Go binary
CMD ["./main"]
