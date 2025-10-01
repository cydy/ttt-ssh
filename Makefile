.PHONY: build run docker-build docker-run docker-stop clean

# Build the Go binary
build:
	go build -o tictactoe-ssh .

# Run locally
run: build
	./tictactoe-ssh

# Build Docker image
docker-build:
	docker-compose build

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker container
docker-stop:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Clean up
clean:
	rm -f tictactoe-ssh
	docker-compose down
