# Use the official Golang image to build your application
FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the Go modules and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the entire application code into the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd/go_auth/main.go

# Expose the port your application uses
EXPOSE 8080

# Command to run the application when the container starts
CMD ["./main"]
