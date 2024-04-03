# Stage 1: Build the application
FROM golang:alpine as builder

# Set the working directory inside the container
WORKDIR /go/src/app

# Enable Go modules
ENV GO111MODULE=on

# Install reflex using go install
RUN go install github.com/cespare/reflex@latest

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o ./run .

# Stage 2: Create the final image
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the executable from the builder stage
COPY --from=builder /go/src/app/run .

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./run"]
