# Use the official Go image as the base
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy all project files into the container
COPY . .

# Download Go module dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o receipt-processor ./cmd/main.go

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./receipt-processor"]
