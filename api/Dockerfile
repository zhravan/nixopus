FROM golang:1.23.6-alpine AS builder
ARG API_PORT
RUN apk add --no-cache make git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./migrations /app/migrations
COPY ../.env .

RUN make build

FROM alpine:3.18
ARG API_PORT
RUN apk add --no-cache ca-certificates bash docker-cli

WORKDIR /app

COPY --from=builder /app/bin/nixopus-api .
COPY --from=builder /app/migrations ./migrations
COPY ../.env .

RUN mkdir -p /etc/nixopus/ssh && \
    chmod 700 /etc/nixopus/ssh && \
    mkdir -p /etc/nixopus/docker-certs && \
    chmod 700 /etc/nixopus/docker-certs && \
    mkdir -p /etc/nixopus/configs && \
    chmod 700 /etc/nixopus/configs

RUN echo "#!/bin/sh" > /app/debug.sh && \
    echo "echo '=== SSH Key Debug Info ==='" >> /app/debug.sh && \
    echo "ls -la /etc/nixopus/ssh/" >> /app/debug.sh && \
    echo "echo '=== SSH Key Content ==='" >> /app/debug.sh && \
    echo "cat /etc/nixopus/ssh/id_rsa" >> /app/debug.sh && \
    echo "echo '=== Environment Variables ==='" >> /app/debug.sh && \
    echo "env | grep SSH" >> /app/debug.sh && \
    chmod +x /app/debug.sh

EXPOSE ${API_PORT}

ENTRYPOINT ["sh", "-c", "/app/debug.sh && /app/nixopus-api"]