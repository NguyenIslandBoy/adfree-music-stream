# Stage 1 — Build Go binary
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 2 — Runtime
FROM python:3.12-alpine

# Install yt-dlp with EJS support
RUN pip install -U "yt-dlp[default]" --quiet

# Install Deno (JS runtime for yt-dlp EJS challenge solver)
RUN apk add --no-cache curl unzip && \
    curl -fsSL https://deno.land/install.sh | sh && \
    ln -s /root/.deno/bin/deno /usr/local/bin/deno

# Install supercronic for cron jobs
RUN wget -qO /usr/local/bin/supercronic \
    https://github.com/aptible/supercronic/releases/download/v0.2.29/supercronic-linux-amd64 && \
    chmod +x /usr/local/bin/supercronic

# Weekly yt-dlp auto-update (every Monday at 3am)
RUN echo "0 3 * * 1 pip install -U 'yt-dlp[default]' --quiet" > /etc/crontab

WORKDIR /app
COPY --from=builder /app/server .
COPY static/ ./static/

EXPOSE 8080

CMD ["sh", "-c", "supercronic /etc/crontab & ./server"]