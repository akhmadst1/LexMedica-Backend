# Use the official Golang image as the base
FROM golang:1.24.1

# Set the working directory in the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Expose the port that the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/app/main"]
