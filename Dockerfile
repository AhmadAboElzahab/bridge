# Use a newer Golang base image
FROM golang:1.24-alpine

# Install system dependencies needed for cgo and libwebp
RUN apk update && apk add --no-cache \
    gcc \
    g++ \
    make \
    libc-dev \
    libwebp-dev \
    libjpeg-turbo-dev \
    zlib-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy the rest of the application code
COPY . .

# Expose the port your app will use
EXPOSE 8080

# Run the app using Air
CMD ["air", "-c", ".air.toml"]
