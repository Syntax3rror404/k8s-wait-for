# Variables
ARG ALPINE_VERSION=3.22
ARG GOLANG_VERSION=1.25.1
ARG USERID=65532

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
ARG USERID

# Create a non-root user "app" to run the application
RUN addgroup -g ${USERID} app && adduser -u ${USERID} -G app -S app

# Copy compiled bin into final image
COPY --from=builder /app/waitfor /usr/local/bin/waitfor

# Set to nonroot user
USER ${USERID}:${USERID}

ENTRYPOINT ["waitfor"]
