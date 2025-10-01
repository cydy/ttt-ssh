# Installation Guide

Step-by-step instructions to install and run SSH Tic-Tac-Toe on various platforms.

## Prerequisites

Choose one installation method:

### Option 1: Docker (Easiest)
- Docker installed
- Docker Compose installed (or use `docker compose`)

### Option 2: Go (For Developers)
- Go 1.21 or higher
- Git (for cloning)

### Option 3: Pre-built Binary
- Download from releases (if available)
- Linux/Mac/Windows compatible

---

## Installation Methods

### Method 1: Docker (Recommended for Production)

#### Step 1: Clone the repository
```bash
git clone <repository-url>
cd tictactoe-ssh
```

#### Step 2: Start with Docker Compose
```bash
docker-compose up -d
```

#### Step 3: Verify it's running
```bash
docker-compose ps
```

#### Step 4: Connect and play!
```bash
ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost
```

#### Stop the server
```bash
docker-compose down
```

---

### Method 2: Build from Source (Go)

#### Step 1: Install Go
If you don't have Go installed:

**Linux/Mac (using snap):**
```bash
sudo snap install go --classic
```

**Mac (using Homebrew):**
```bash
brew install go
```

**Or download from:** https://golang.org/dl/

#### Step 2: Clone the repository
```bash
git clone <repository-url>
cd tictactoe-ssh
```

#### Step 3: Download dependencies
```bash
go mod download
```

#### Step 4: Build the binary
```bash
go build -o tictactoe-ssh .
```

#### Step 5: Run the server
```bash
./tictactoe-ssh
```

Or use the quick start script:
```bash
chmod +x run.sh
./run.sh
```

#### Step 6: Connect and play!
```bash
ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost
```

---

### Method 3: Direct Run (No Build)

```bash
# Clone the repository
git clone <repository-url>
cd tictactoe-ssh

# Run directly
go run main.go
```

---

## Platform-Specific Instructions

### Ubuntu/Debian

```bash
# Install dependencies
sudo apt update
sudo apt install -y git golang-go

# Clone and build
git clone <repository-url>
cd tictactoe-ssh
go build -o tictactoe-ssh .

# Run
./tictactoe-ssh
```

### CentOS/RHEL

```bash
# Install dependencies
sudo yum install -y git golang

# Clone and build
git clone <repository-url>
cd tictactoe-ssh
go build -o tictactoe-ssh .

# Run
./tictactoe-ssh
```

### macOS

```bash
# Install Homebrew if not installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Go
brew install go

# Clone and build
git clone <repository-url>
cd tictactoe-ssh
go build -o tictactoe-ssh .

# Run
./tictactoe-ssh
```

### Windows

#### Option A: Using WSL2 (Recommended)
```powershell
# Install WSL2 and Ubuntu
wsl --install

# Then follow Ubuntu instructions above
```

#### Option B: Native Windows
```powershell
# Install Go from https://golang.org/dl/
# Then in PowerShell:

git clone <repository-url>
cd tictactoe-ssh
go build -o tictactoe-ssh.exe .
.\tictactoe-ssh.exe
```

Connect using:
```powershell
ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost
```

Or use PuTTY:
- Host: localhost
- Port: 2222
- Connection type: SSH

---

## VPS/Server Installation

### DigitalOcean Droplet

```bash
# SSH into your droplet
ssh root@your-droplet-ip

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Clone and run
git clone <repository-url>
cd tictactoe-ssh
docker-compose up -d

# Configure firewall
sudo ufw allow 2222/tcp
sudo ufw allow 22/tcp  # Keep SSH access!
sudo ufw enable
```

Connect from anywhere:
```bash
ssh -p 2222 your-droplet-ip
```

### AWS EC2

```bash
# Launch EC2 instance (Ubuntu 22.04 recommended)
# Configure Security Group to allow port 2222

# SSH into instance
ssh -i your-key.pem ubuntu@your-ec2-ip

# Install Docker
sudo apt update
sudo apt install -y docker.io docker-compose
sudo usermod -aG docker ubuntu

# Clone and run
git clone <repository-url>
cd tictactoe-ssh
docker-compose up -d
```

### Google Cloud Platform

```bash
# Create VM instance
# Configure firewall rule for port 2222

# SSH into instance
gcloud compute ssh your-instance-name

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Clone and run
git clone <repository-url>
cd tictactoe-ssh
docker-compose up -d
```

---

## Configuration

### Change Port

#### Method 1: Environment Variable
```bash
export SSH_PORT=3000
./tictactoe-ssh
```

#### Method 2: Docker Compose
Edit `docker-compose.yml`:
```yaml
environment:
  - SSH_PORT=3000
ports:
  - "3000:3000"
```

### Run as System Service (Linux)

Create `/etc/systemd/system/tictactoe.service`:
```ini
[Unit]
Description=SSH Tic-Tac-Toe Server
After=network.target

[Service]
Type=simple
User=tictactoe
WorkingDirectory=/opt/tictactoe-ssh
ExecStart=/opt/tictactoe-ssh/tictactoe-ssh
Restart=always
Environment="SSH_PORT=2222"

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable tictactoe
sudo systemctl start tictactoe
sudo systemctl status tictactoe
```

---

## Verification

### Check if server is running

```bash
# Check process
ps aux | grep tictactoe

# Check port
netstat -tulpn | grep 2222
# or
lsof -i :2222

# Test connection
ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost
```

### View logs

#### Docker:
```bash
docker-compose logs -f
```

#### Systemd:
```bash
sudo journalctl -u tictactoe -f
```

---

## Troubleshooting

### Port already in use
```bash
# Find what's using the port
lsof -i :2222

# Kill the process or change port
export SSH_PORT=2223
```

### Permission denied
```bash
# Docker permission
sudo usermod -aG docker $USER
# Log out and back in

# File permission
chmod +x tictactoe-ssh
chmod +x run.sh
```

### Module download failed
```bash
# Clear module cache
go clean -modcache

# Download again
go mod download
```

### Docker compose not found
```bash
# Try with space (newer syntax)
docker compose up -d

# Or install docker-compose
sudo apt install docker-compose
```

### Connection refused
```bash
# Check if server is running
docker-compose ps
# or
ps aux | grep tictactoe

# Check firewall
sudo ufw status
sudo ufw allow 2222/tcp
```

---

## Uninstallation

### Docker
```bash
docker-compose down
docker rmi tictactoe-ssh_tictactoe
```

### Binary
```bash
rm tictactoe-ssh
```

### System Service
```bash
sudo systemctl stop tictactoe
sudo systemctl disable tictactoe
sudo rm /etc/systemd/system/tictactoe.service
sudo systemctl daemon-reload
```

---

## Next Steps

After installation:

1. Read [QUICKSTART.md](QUICKSTART.md) for gameplay
2. Check [EXAMPLES.md](EXAMPLES.md) for usage examples
3. Review [README.md](README.md) for full documentation
4. See [ARCHITECTURE.md](ARCHITECTURE.md) for technical details

Enjoy playing! 🎮
