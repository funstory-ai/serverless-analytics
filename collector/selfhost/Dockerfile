# Build step
FROM golang:1.23 AS builder
RUN mkdir -p /build
WORKDIR /build
COPY . .
RUN go build -trimpath -ldflags "-s -w" -o bin/collect ./cmd

# Final step
FROM debian:12-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/* \
EXPOSE 8080
WORKDIR /app
COPY --from=builder /build/bin/collect /app/collect
ENTRYPOINT ["/app/collect"]