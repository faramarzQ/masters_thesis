# syntax=docker/dockerfile:1

FROM golang:1.20-alpine

RUN apk add make

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
# RUN go mod download

# Copy the source code.
COPY . ./

# Build 
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/scaler ./cmd/scaler/main.go

cmd ["sleep", "500"]