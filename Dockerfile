# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN apk add --no-cache git make && \
    make build

# Stage 2: Minimal runtime
# No ffmpeg: Instagram serves single-file MP4s and the bot requests only
# pre-muxed formats (see downloadInstagramVideo), so no stream merging is needed.
# yt-dlp is pure Python, so python3 is the only runtime dependency.
FROM alpine:3.21

RUN apk add --no-cache python3 py3-pip ca-certificates && \
    pip3 install --no-cache-dir --break-system-packages yt-dlp && \
    apk del py3-pip && \
    rm -rf /root/.cache && \
    yt-dlp --version

WORKDIR /app

COPY --from=builder /go/src/app/kbot .

ENTRYPOINT ["./kbot", "start"]
