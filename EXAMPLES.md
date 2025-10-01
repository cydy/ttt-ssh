# Game Examples

This document provides examples of how to play SSH Tic-Tac-Toe.

## Example 1: Playing Against CPU

```
$ ssh -p 2222 localhost

╔════════════════════════════════════╗
║   Welcome to SSH Tic-Tac-Toe!     ║
╚════════════════════════════════════╝

Choose game mode:
1. Play against CPU
2. Multiplayer (Create room)
3. Multiplayer (Join room)
4. Quit

Enter choice (1-4): 1

=== Playing against CPU ===
You are: X
CPU is: O

Enter moves as: row col (e.g., '1 2' for row 1, column 2)

     1   2   3
   ┌───┬───┬───┐
 1 │   │   │   │
   ├───┼───┼───┤
 2 │   │   │   │
   ├───┼───┼───┤
 3 │   │   │   │
   └───┴───┴───┘

Your turn (X): 1 1

     1   2   3
   ┌───┬───┬───┐
 1 │ X │   │   │
   ├───┼───┼───┤
 2 │   │   │   │
   ├───┼───┼───┤
 3 │   │   │   │
   └───┴───┴───┘

CPU played: 2 2

     1   2   3
   ┌───┬───┬───┐
 1 │ X │   │   │
   ├───┼───┼───┤
 2 │   │ O │   │
   ├───┼───┼───┤
 3 │   │   │   │
   └───┴───┴───┘

Your turn (X): 1 2
...
```

## Example 2: Creating a Multiplayer Room

**Player 1:**
```
$ ssh -p 2222 localhost

Choose game mode:
1. Play against CPU
2. Multiplayer (Create room)
3. Multiplayer (Join room)
4. Quit

Enter choice (1-4): 2

=== Multiplayer Room Created ===
Room Code: a3b7f9
You are: X
Waiting for opponent to join...

Opponent joined!

     1   2   3
   ┌───┬───┬───┐
 1 │   │   │   │
   ├───┼───┼───┤
 2 │   │   │   │
   ├───┼───┼───┤
 3 │   │   │   │
   └───┴───┴───┘

Room: a3b7f9
Players: X O
Current turn: X

Your turn (X): 2 2
```

**Player 2:**
```
$ ssh -p 2222 localhost

Choose game mode:
1. Play against CPU
2. Multiplayer (Create room)
3. Multiplayer (Join room)
4. Quit

Enter choice (1-4): 3

Enter room code: a3b7f9

=== Joined Room a3b7f9 ===
You are: O

     1   2   3
   ┌───┬───┬───┐
 1 │   │   │   │
   ├───┼───┼───┤
 2 │   │ X │   │
   ├───┼───┼───┤
 3 │   │   │   │
   └───┴───┴───┘

Room: a3b7f9
Players: X O
Current turn: O

Your turn (O): 1 1
```

## Example 3: Winning Game

```
     1   2   3
   ┌───┬───┬───┐
 1 │ X │ X │ X │
   ├───┼───┼───┤
 2 │ O │ O │   │
   ├───┼───┼───┤
 3 │   │   │   │
   └───┴───┴───┘

🎉 Congratulations! You won!

Press Enter to return to main menu...
```

## Example 4: Draw Game

```
     1   2   3
   ┌───┬───┬───┐
 1 │ X │ O │ X │
   ├───┼───┼───┤
 2 │ O │ X │ X │
   ├───┼───┼───┤
 3 │ O │ X │ O │
   └───┴───┴───┘

🤝 It's a draw!

Press Enter to return to main menu...
```

## Tips for Playing

1. **Board Coordinates**: The board uses a coordinate system where:
   - First number is the row (1-3 from top to bottom)
   - Second number is the column (1-3 from left to right)

2. **CPU Strategy**: The AI follows this priority:
   - Win if possible
   - Block opponent from winning
   - Take center (2 2)
   - Take corners
   - Take any available position

3. **Multiplayer Tips**:
   - Room codes are 6 characters long
   - Share the code with your friend before they try to join
   - Only 2 players per room
   - Rooms are automatically deleted after the game ends

4. **Common Commands**:
   - Enter moves: `1 2` (row 1, column 2)
   - Invalid moves show an error - try again
   - Press Enter to continue after game ends

## Deployment Examples

### Running Locally

```bash
# Build and run
go build -o tictactoe-ssh .
./tictactoe-ssh
```

### Running with Docker

```bash
# Using docker-compose
docker-compose up -d

# Or build and run manually
docker build -t tictactoe-ssh .
docker run -d -p 2222:2222 --name tictactoe tictactoe-ssh
```

### Running on a Server

```bash
# On your server
cd /opt/tictactoe-ssh
docker-compose up -d

# Configure firewall
sudo ufw allow 2222/tcp

# Connect from anywhere
ssh -p 2222 your-server.com
```

### Custom Port

```bash
# Set environment variable
export SSH_PORT=3000
./tictactoe-ssh

# Or with Docker
docker run -d -p 3000:3000 -e SSH_PORT=3000 tictactoe-ssh
```
