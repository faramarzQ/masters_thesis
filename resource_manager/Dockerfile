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
RUN make build-scaler

cmd ["sleep", "500"]