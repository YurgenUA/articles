# Use the official Golang image as the base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Copy the assets folder into the container
COPY assets ./assets

# Expose the port the application runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]