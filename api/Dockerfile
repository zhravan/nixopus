FROM golang:1.21-alpine AS builder

ARG PORT=8443
ENV PORT=${PORT}

RUN apk add --no-cache make git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./migrations /app/migrations

RUN make build

FROM alpine:3.18

ARG PORT=8443
ENV PORT=${PORT}

RUN apk add --no-cache ca-certificates bash docker-cli

WORKDIR /app

COPY --from=builder /app/bin/nixopus-api .
COPY --from=builder /app/migrations ./migrations

EXPOSE ${PORT}

ENTRYPOINT ["./nixopus-api"]