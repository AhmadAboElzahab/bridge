# Use a newer Golang base image (1.24 or later)
FROM golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy the rest of the application code
COPY . .

# Expose the port your app runs on
EXPOSE 8080

# Set the default command to run Air with the correct entry point
CMD ["air", "-c", ".air.toml"]

