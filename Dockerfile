# Builder-Stage
FROM golang:1.25.1-alpine3.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o waitfor

# Runner
FROM alpine:3.22.1

# Create a non-root user "app" to run the application
RUN addgroup -g 1000 app && adduser -u 1000 -G app -S app

COPY --from=builder /app/waitfor /usr/local/bin/waitfor

USER 1000:1000

ENTRYPOINT ["waitfor"]