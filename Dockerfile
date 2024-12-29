FROM golang:1.23-alpine

# Install air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Command to run air
CMD ["air"]