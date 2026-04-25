# Stage 1 — Build Go binary
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 2 — Runtime
FROM python:3.12-alpine

# Install yt-dlp
RUN pip install -U yt-dlp --quiet

# Install supercronic for cron jobs (Alpine has no cron by default)
RUN wget -qO /usr/local/bin/supercronic \
    https://github.com/aptible/supercronic/releases/download/v0.2.29/supercronic-linux-amd64 && \
    chmod +x /usr/local/bin/supercronic

# Weekly yt-dlp auto-update (every Monday at 3am)
RUN echo "0 3 * * 1 pip install -U yt-dlp --quiet" > /etc/crontab

WORKDIR /app
COPY --from=builder /app/server .
COPY static/ ./static/

EXPOSE 8080

# Fix: JSON format for CMD
CMD ["sh", "-c", "supercronic /etc/crontab & ./server"]