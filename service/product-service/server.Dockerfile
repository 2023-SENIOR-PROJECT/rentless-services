# Use an official Golang runtime as a parent image
FROM golang:1.21 as builder

# Set the working directory
WORKDIR /go/src

# Copy the local package files to the container's workspace
COPY . .
EXPOSE 50051
# EXPOSE 8081
# Build the Go app
RUN go build -o main .
CMD ["./main"]
# RUN go build grpc-client/client.go

# Run the binary of ./main and ./client
# FROM alpine:latest
# WORKDIR /root/
# COPY --from=builder /go/src/main .
# CMD ["./main"]
# RUN apk add --no-cache bash
# WORKDIR /root/
# COPY --from=builder /go/src/main .
# COPY --from=builder /go/src/client .
# COPY --from=builder /go/src/start.sh .
# # RUN ls -al
# RUN chmod +x /root/start.sh
# CMD ["/root/start.sh"]
