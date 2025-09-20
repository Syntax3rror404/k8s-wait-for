# Variables
ARG ALPINE_VERSION=3.22
ARG GOLANG_VERSION=1.25.1

FROM alpine:${ALPINE_VERSION} AS builder

#
# Builder-Stage ++++++++++++++++++++++++++++++++++++++++
#
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o waitfor

#
# Final-Stage ++++++++++++++++++++++++++++++++++++++++
#
FROM alpine:${ALPINE_VERSION}

# Create a non-root user "app" to run the application
RUN addgroup -g 65532 app && adduser -u 65532 -G app -S app

COPY --from=builder /app/waitfor /usr/local/bin/waitfor

USER 65532:65532

ENTRYPOINT ["waitfor"]
