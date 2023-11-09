# Use an official Golang runtime as a parent image
FROM golang:1.21 as builder

# Set the working directory
WORKDIR /go/src

# Copy the local package files to the container's workspace
COPY . .
EXPOSE 8083
# Build the Go app
RUN go build -o main .
CMD ["./main"]