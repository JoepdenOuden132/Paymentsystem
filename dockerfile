# Use the official Golang image as the builder stage
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code
COPY . ./

# Build the Go application
RUN go build -o /restapi

# Use a minimal image for the final runtime environment
FROM alpine:latest

# Set the working directory in the runtime container
WORKDIR /

# Copy the built binary from the builder image
COPY --from=builder /restapi /restapi

# Expose the port on which the application will run
EXPOSE 8080

# Define the command to run the application
CMD ["/restapi"]