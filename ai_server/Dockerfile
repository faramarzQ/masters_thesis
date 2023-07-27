# syntax=docker/dockerfile:1

FROM python:3.9-alpine

# Set destination for COPY
WORKDIR /app

# Install modules
RUN pip install numpy

# Copy the source code.
COPY ./ ./

EXPOSE 8080

# Run
CMD ["python3", "./main.py"] 