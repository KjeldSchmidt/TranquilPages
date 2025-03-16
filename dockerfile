FROM golang:1.24-alpine AS build
RUN apk add --no-cache gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Enable CGO for SQLite
ENV CGO_ENABLED=1

# Build the application
RUN go build -o tranquil-pages

# Use a lightweight image for runtime
FROM alpine:latest

# Install SQLite dependency
RUN apk add --no-cache sqlite-libs

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=build /app/tranquil-pages .

# Expose the application's port
EXPOSE 8080

# Run the application
CMD ["./tranquil-pages"]
