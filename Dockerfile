FROM golang:1.23-alpine

# # Install air
# RUN go install github.com/air-verse/air@latest

# WORKDIR /app

# # Copy go.mod and go.sum files
# COPY go.mod go.sum ./

# # Download dependencies
# RUN go mod download

# # Copy the rest of the application code
# COPY . .

# # Command to run air
# CMD ["air"]
# syntax=docker/dockerfile:1

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /location-share-backend

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8000

# Run
CMD ["/location-share-backend"]
