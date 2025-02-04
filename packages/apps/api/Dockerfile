FROM golang:1.23-bookworm as builder

WORKDIR /app/api

# Install dependencies
COPY ./api/go.* /app/api/
COPY ./common/go.* /app/common/
COPY ./db/go.* /app/db/

RUN go mod download

# Build
COPY ./api /app/api
COPY ./common /app/common
COPY ./db /app/db/

RUN go build -v -o /server ./cmd/run_server

# Create a minimal image
FROM debian:bookworm-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /server /server

CMD ["/server"]
