# Use the official Golang image as the builder stage
# Use the official Golang image as the builder stage
FROM golang:1.22.1-bullseye AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code
COPY . ./

# Build the Go application and verify the binary is created
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /restapi && ls -lah /restapi

# Use a minimal image for the final runtime environment
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Set the working directory in the runtime container
WORKDIR /app

# Copy the built binary from the builder image
COPY --from=builder /restapi /app/restapi

# Expose the port on which the application will run
EXPOSE 80

# Define the command to run the application
CMD ["/app/restapi"]