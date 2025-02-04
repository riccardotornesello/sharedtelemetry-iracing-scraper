FROM golang:1.23-bookworm as builder

WORKDIR /app/drivers

# Install dependencies
COPY ./drivers/go.* /app/drivers/
COPY ./common/go.* /app/common/
COPY ./db/go.* /app/db/
COPY ./iracing-api/go.* /app/iracing-api/

RUN go mod download

# Build
COPY ./drivers /app/drivers
COPY ./common /app/common
COPY ./db /app/db
COPY ./iracing-api /app/iracing-api/

ARG COMMAND

RUN go build -v -o /server ./cmd/${COMMAND}

# Create a minimal image
FROM debian:bookworm-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /server /server

CMD ["/server"]
