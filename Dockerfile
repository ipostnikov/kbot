
FROM quay.io/projectquay/golang:1.24 AS builder
WORKDIR /go/src/app
COPY . .
RUN make build

FROM alpine:latest
WORKDIR /
RUN apk add --no-cache python3 py3-pip ca-certificates curl && \
    pip install --no-cache-dir yt-dlp
COPY --from=builder /go/src/app/kbot .
ENTRYPOINT [ "./kbot", "start" ]