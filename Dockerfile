# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN apk add --no-cache git make bash curl && \
    make build && \
    apk del git make bash curl

# Stage 2: Final image with yt-dlp and ffmpeg
FROM alpine:3.19

RUN apk add --no-cache python3 ffmpeg curl ca-certificates && \
    curl -fL https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod +x /usr/local/bin/yt-dlp && \
    yt-dlp --version

WORKDIR /app

# Copy compiled Go binary
COPY --from=builder /go/src/app/kbot .

ENTRYPOINT ["./kbot", "start"]
