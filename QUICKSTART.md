# Quick Start Guide

Get up and running with SSH Tic-Tac-Toe in minutes!

## 🚀 Fastest Way to Start

### Option 1: Using Docker (Recommended)

```bash
# Start the server
docker-compose up -d

# Connect and play
ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost
```

That's it! The game is now running on port 2222.

### Option 2: Using Go

```bash
# Build and run
go build -o tictactoe-ssh .
./tictactoe-ssh

# Or use the convenience script
./run.sh
```

### Option 3: Direct Run

```bash
go run main.go
```

## 🎮 How to Play

### 1. Connect via SSH

```bash
ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost
```

### 2. Choose Your Game Mode

You'll see a menu:
```
1. Play against CPU
2. Multiplayer (Create room)
3. Multiplayer (Join room)
4. Quit
```

### 3. Make Moves

Enter moves as two numbers: `row column`

Example: `1 2` means row 1, column 2

```
     1   2   3
   ┌───┬───┬───┐
 1 │   │ X │   │  ← Your move at 1 2
   ├───┼───┼───┤
 2 │   │   │   │
   ├───┼───┼───┤
 3 │   │   │   │
   └───┴───┴───┘
```

## 🤖 Playing Against CPU

1. Choose option `1`
2. You play as `X`, CPU plays as `O`
3. Enter your moves when prompted
4. First to get 3 in a row wins!

## 👥 Multiplayer Mode

### Create a Room (Player 1)

1. Choose option `2`
2. You'll get a 6-character room code (e.g., `a3b7f9`)
3. Share this code with your friend
4. Wait for them to join

### Join a Room (Player 2)

1. Choose option `3`
2. Enter the room code your friend shared
3. Game starts immediately!

## 🛠️ Configuration

### Change Port

```bash
# Set environment variable
export SSH_PORT=3000
./tictactoe-ssh

# Or with Docker
docker run -p 3000:3000 -e SSH_PORT=3000 tictactoe-ssh
```

### Run on Server

```bash
# On your VPS/server
git clone <your-repo>
cd tictactoe-ssh
docker-compose up -d

# Allow port through firewall
sudo ufw allow 2222/tcp

# Connect from anywhere
ssh -p 2222 your-server.com
```

## 📋 Useful Commands

```bash
# Build the binary
go build -o tictactoe-ssh .

# Run tests
go test ./...

# Build Docker image
docker-compose build

# View Docker logs
docker-compose logs -f

# Stop Docker container
docker-compose down

# Clean up
make clean
```

## 🐛 Troubleshooting

### Can't connect?
- Check if server is running: `docker-compose ps` or `ps aux | grep tictactoe`
- Verify port 2222 is not in use: `lsof -i :2222`
- Try with verbose SSH: `ssh -v -p 2222 localhost`

### Room not found?
- Make sure the room creator hasn't disconnected
- Check that you typed the code correctly (case-sensitive)
- Create a new room if needed

### Build errors?
- Ensure Go 1.21+ is installed: `go version`
- Download dependencies: `go mod download`
- Clean and rebuild: `go clean && go build`

## 📚 Next Steps

- Read [README.md](README.md) for detailed documentation
- Check [EXAMPLES.md](EXAMPLES.md) for game examples
- View the code in `pkg/` to understand how it works
- Customize the game by modifying the source

## 🎯 Tips

1. **CPU Strategy**: The AI uses balanced difficulty - sometimes strategic, sometimes random for easier gameplay
2. **Coordinates**: Remember it's `row col` not `x y`
3. **Multiplayer**: First player is always X, second player is always O
4. **Room Codes**: They're short (6 chars) for easy sharing

Happy gaming! 🎮
