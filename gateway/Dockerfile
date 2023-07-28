# syntax=docker/dockerfile:1

FROM golang:1.20-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway

EXPOSE 4444

# Run
CMD ["/gateway"]
