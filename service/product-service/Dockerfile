# Start from the official golang image
FROM golang:1.21.4-alpine3.17

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files from the parent directory
COPY ../../go.mod ./go.mod
COPY ../../go.sum ./go.sum
# RUN go mod download

# Copy the Go application source code from the current directory
COPY . .

# Build the Go application
RUN go build -o product-service

# Expose the port the service will be running on
EXPOSE 50051

# Command to run the executable
CMD ["./product-service"]
