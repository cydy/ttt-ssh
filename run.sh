#!/bin/bash

# SSH Tic-Tac-Toe Quick Start Script

echo "🎮 SSH Tic-Tac-Toe Server"
echo "=========================="
echo ""

# Check if binary exists
if [ ! -f "./tictactoe-ssh" ]; then
    echo "Building application..."
    go build -o tictactoe-ssh .
    if [ $? -ne 0 ]; then
        echo "❌ Build failed!"
        exit 1
    fi
    echo "✅ Build successful!"
    echo ""
fi

# Set default port if not set
export SSH_PORT=${SSH_PORT:-2222}

echo "Starting server on port $SSH_PORT..."
echo ""
echo "📡 Connect with: ssh -p $SSH_PORT localhost"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

./tictactoe-ssh
