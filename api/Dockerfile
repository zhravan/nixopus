FROM golang:1.21-alpine AS builder
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

EXPOSE ${API_PORT}

ENTRYPOINT ["./nixopus-api"]