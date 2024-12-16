FROM golang:1.23-bookworm as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

ARG PACKAGE

RUN go build -v -o server ./cmd/${PACKAGE}

FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /server

CMD ["/server"]
