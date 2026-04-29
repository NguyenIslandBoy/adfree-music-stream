.PHONY: build test vet lint clean docker-build docker-up docker-down update-ytdlp

# Build the Go binary
build:
	go build -o server ./cmd/server/

# Run all tests
test:
# # 	go test ./... -v

# Run Go vet (static analysis)
# # vet:
	go vet ./...

# # Run build + vet + test in sequence
check: vet test

# # # # Clean built binary
# clean:
# 	rm -f server

# Docker
docker-build:
	docker build -t adfree-music-stream .

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# Update yt-dlp
update-ytdlp:
	pip install -U yt-dlp

# Run locally (no Docker)
run: build
	./server