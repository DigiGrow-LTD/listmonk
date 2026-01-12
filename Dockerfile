# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk --no-cache add make git nodejs npm yarn

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy frontend package files and install deps (cached layer)
COPY frontend/package.json frontend/yarn.lock frontend/
RUN cd frontend && yarn install --frozen-lockfile

# Copy email-builder package files and install deps (cached layer)
COPY frontend/email-builder/package.json frontend/email-builder/yarn.lock frontend/email-builder/
RUN cd frontend/email-builder && yarn install --frozen-lockfile

# Copy rest of source code
COPY . .

# Build the application (deps already installed, just build)
RUN make dist

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata shadow su-exec

WORKDIR /listmonk

COPY --from=builder /app/listmonk .
COPY config.toml.sample config.toml
COPY docker-entrypoint.sh /usr/local/bin/

RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 9000

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["./listmonk"]
