# SSH Tic-Tac-Toe

A text-based Tic-Tac-Toe game accessible via SSH, built with Go. Features both CPU opponent and multiplayer modes with room codes.

## Features

- 🎮 **Text-based UI** - Beautiful ASCII art game board
- 🤖 **CPU Opponent** - Play against an AI with balanced difficulty
- 👥 **Multiplayer Mode** - Play with friends using room codes
- 🔐 **SSH Access** - Connect from anywhere using SSH
- 🐳 **Docker Support** - Easy deployment with Docker

## Quick Start with Docker

The easiest way to run the game is using Docker:

```bash
# Build and run with docker-compose
docker-compose up -d

# Or using make
make docker-run
```

Connect to the game:
```bash
ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost
```

## Manual Setup

### Prerequisites

- Go 1.21 or higher
- (Optional) Docker and Docker Compose

### Build from Source

```bash
# Clone the repository
git clone <your-repo-url>
cd tictactoe-ssh

# Download dependencies
go mod download

# Build the binary
go build -o tictactoe-ssh .

# Run the server
./tictactoe-ssh
```

## Usage

### Connecting to the Game

```bash
# Connect to localhost (default port 2222)
ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost

# Connect to a remote server
ssh -p 2222 username@your-server.com
```

### Game Modes

#### 1. Play Against CPU
- Choose option 1 from the main menu
- You play as X, CPU plays as O
- Enter moves as `row col` (e.g., `1 2` for row 1, column 2)

#### 2. Multiplayer - Create Room
- Choose option 2 from the main menu
- You'll receive a 6-character room code
- Share this code with your friend
- Wait for them to join
- You play as X (first player)

#### 3. Multiplayer - Join Room
- Choose option 3 from the main menu
- Enter the room code provided by your friend
- You play as O (second player)

### Game Controls

- Enter moves as: `row col` (e.g., `1 2`)
- Board coordinates:
  ```
       1   2   3
     ┌───┬───┬───┐
   1 │   │   │   │
     ├───┼───┼───┤
   2 │   │   │   │
     ├───┼───┼───┤
   3 │   │   │   │
     └───┴───┴───┘
  ```

## Docker Commands

```bash
# Build the Docker image
make docker-build

# Run the container
make docker-run

# View logs
make logs

# Stop the container
make docker-stop

# Clean up
make clean
```

## Configuration

The SSH port can be configured via environment variable:

```bash
# Set custom port
export SSH_PORT=3000
./tictactoe-ssh
```

Or in docker-compose.yml:
```yaml
environment:
  - SSH_PORT=3000
ports:
  - "3000:3000"
```

## Project Structure

```
tictactoe-ssh/
├── main.go                 # Application entry point
├── pkg/
│   ├── game/
│   │   ├── board.go       # Game board logic
│   │   ├── ai.go          # CPU opponent AI
│   │   └── multiplayer.go # Multiplayer room management
│   └── server/
│       └── server.go      # SSH server implementation
├── Dockerfile             # Docker build configuration
├── docker-compose.yml     # Docker Compose configuration
├── Makefile              # Build and run commands
├── go.mod                # Go module definition
└── README.md             # This file
```

## How It Works

### SSH Server
- Uses `gliderlabs/ssh` library for SSH server functionality
- No authentication required (public access)
- Each connection gets its own game session

### CPU AI
- Balanced difficulty algorithm:
  1. Try to win (50% of the time)
  2. Block opponent from winning (30% of the time)
  3. Take center if available (40% of the time)
  4. Take corners (30% of the time)
  5. Make random moves (more often for easier gameplay)

### Multiplayer
- Room-based system with 6-character codes
- First player creates room (plays as X)
- Second player joins with code (plays as O)
- Rooms are automatically cleaned up after game ends

## Deployment

### Deploy to a VPS

1. **Copy files to server:**
   ```bash
   scp -r . user@your-server:/path/to/tictactoe-ssh
   ```

2. **SSH into server:**
   ```bash
   ssh user@your-server
   cd /path/to/tictactoe-ssh
   ```

3. **Run with Docker:**
   ```bash
   docker-compose up -d
   ```

4. **Configure firewall (if needed):**
   ```bash
   sudo ufw allow 2222/tcp
   ```

### Deploy to Cloud

The application can be deployed to any cloud provider that supports Docker:

- **AWS ECS/Fargate**
- **Google Cloud Run**
- **Azure Container Instances**
- **DigitalOcean App Platform**
- **Heroku**

Simply build the Docker image and deploy it with port 2222 exposed.

## Development

### Run locally without Docker:

```bash
# Install dependencies
go mod download

# Run the server
go run main.go

# Or use the convenience script
./run.sh
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests in verbose mode
go test ./... -v
```

### Run with hot reload (using air):

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

## Troubleshooting

### "Connection refused"
- Ensure the server is running: `docker-compose ps`
- Check if port 2222 is available: `lsof -i :2222`
- Verify firewall settings

### "Room not found"
- Room codes are case-sensitive
- Rooms are deleted after games end
- Ensure the room creator hasn't disconnected

### SSH connection issues
- Try with verbose mode: `ssh -v -p 2222 localhost`
- Ensure SSH client is installed
- Check if port is accessible from your network

## License

MIT License - feel free to use and modify as needed.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
