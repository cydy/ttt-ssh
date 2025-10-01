# Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     SSH Clients                         │
│  (ssh -p 2222 localhost from multiple terminals)       │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                  SSH Server (port 2222)                 │
│              (gliderlabs/ssh library)                   │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                  Session Handler                        │
│            (pkg/server/server.go)                       │
│                                                          │
│  ┌──────────────────────────────────────────────┐      │
│  │  Game Mode Selection Menu                    │      │
│  │  1. CPU Mode                                  │      │
│  │  2. Create Multiplayer Room                  │      │
│  │  3. Join Multiplayer Room                    │      │
│  └──────────────────────────────────────────────┘      │
└─────────┬───────────────────┬────────────────────┬──────┘
          │                   │                    │
          ▼                   ▼                    ▼
    ┌─────────┐         ┌──────────┐        ┌──────────┐
    │CPU Mode │         │Create    │        │Join      │
    │         │         │Room      │        │Room      │
    └────┬────┘         └────┬─────┘        └────┬─────┘
         │                   │                    │
         │                   └────────┬───────────┘
         │                            │
         ▼                            ▼
    ┌─────────────┐          ┌───────────────────┐
    │ Game Logic  │          │  Room Manager     │
    │ (pkg/game)  │◄─────────│  (GlobalRM)       │
    └─────────────┘          └───────────────────┘
```

## Component Breakdown

### 1. SSH Server Layer
**File**: `main.go`, `pkg/server/server.go`
- Listens on port 2222 (configurable)
- Creates session for each connection
- No authentication (public access)
- Handles multiple concurrent connections

### 2. Session Handler
**File**: `pkg/server/server.go`
- Main menu loop
- Route to game modes
- Input/output formatting
- User interaction flow

### 3. Game Logic
**File**: `pkg/game/board.go`
- 3x3 board representation
- Move validation
- Win detection (rows, cols, diagonals)
- Draw detection
- Available moves calculation

### 4. CPU AI
**File**: `pkg/game/ai.go`
```
Decision Tree:
1. Can I win? → Make winning move
2. Can opponent win? → Block them
3. Is center free? → Take center
4. Are corners free? → Take random corner
5. Otherwise → Take any available space
```

### 5. Multiplayer System
**File**: `pkg/game/multiplayer.go`
- Global room manager (singleton)
- Room creation with UUID-based codes
- Player assignment (X/O)
- Turn management
- Thread-safe with mutexes

## Data Flow

### CPU Game Flow
```
Player Connects
      ↓
Choose CPU Mode
      ↓
Create Board
      ↓
┌─────────────────┐
│ Game Loop:      │
│  1. Display     │
│  2. Check Win   │
│  3. Get Move    │──→ Invalid? → Loop back
│  4. Apply Move  │
│  5. Check Win   │──→ Won? → End game
│  6. CPU Move    │
│  7. Check Win   │──→ Won? → End game
└─────────────────┘
      ↓
Return to Menu
```

### Multiplayer Game Flow
```
Player 1 Creates Room          Player 2 Joins Room
        ↓                              ↓
Generate Room Code  ─────Share Code───→ Enter Code
        ↓                              ↓
Wait for P2 ←─────────Connect─────────┘
        ↓
┌──────────────────────────────────────────┐
│ Game Loop (Synchronized):                │
│                                           │
│  P1 Turn:                                 │
│   - P1 makes move                         │
│   - P2 polls for board change             │
│                                           │
│  P2 Turn:                                 │
│   - P2 makes move                         │
│   - P1 polls for board change             │
│                                           │
│  Check win/draw after each move           │
└──────────────────────────────────────────┘
        ↓
Clean up room
```

## Thread Safety

### Multiplayer Rooms
```go
type Room struct {
    mutex sync.Mutex  // Protects all room state
    // ... other fields
}

// All operations lock before modifying:
func (r *Room) MakeMove() {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    // ... safe operations
}
```

### Room Manager
```go
type RoomManager struct {
    rooms map[string]*Room
    mutex sync.RWMutex  // Read-write lock
}

// Creates use write lock
// Joins use read lock
```

## State Management

### Board State
```
cells: [9]Cell  // Array of 9 cells
                // 0-2: Row 1
                // 3-5: Row 2  
                // 6-8: Row 3

Cell values:
- Empty (0)
- X (1)
- O (2)
```

### Room State
```
Room {
    ID: string           // 6-char code
    Board: *Board        // Game board
    Player1: *Player     // X
    Player2: *Player     // O
    Current: Cell        // Whose turn
    GameOver: bool       // Game finished
    Winner: Cell         // Winner or Empty
}
```

## Network Protocol

### SSH Connection
```
Client                          Server
  │                               │
  ├──── TCP Connect ─────────────→│
  │                               │
  ├──── SSH Handshake ───────────→│
  │←──── SSH Response ────────────┤
  │                               │
  │                               │
  │←──── Welcome Message ─────────┤
  │←──── Menu Display ────────────┤
  │                               │
  ├──── Menu Choice ─────────────→│
  │                               │
  │←──── Game Board ──────────────┤
  │                               │
  ├──── Move Input ──────────────→│
  │                               │
  │←──── Updated Board ───────────┤
  │                               │
  │      (repeat until game ends)  │
  │                               │
  ├──── Disconnect ──────────────→│
  │                               │
```

## Deployment Architecture

### Docker Container
```
┌─────────────────────────────────────┐
│  Alpine Linux Container             │
│                                     │
│  ┌───────────────────────────┐    │
│  │  tictactoe-ssh binary     │    │
│  │  (static, 6MB)             │    │
│  └───────────────────────────┘    │
│                                     │
│  Exposed Port: 2222                │
└─────────────────────────────────────┘
```

### Multi-Container (Future)
```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│ Load         │    │ Game         │    │ Redis        │
│ Balancer     │───→│ Server       │───→│ (Rooms)      │
│              │    │ (Multiple)   │    │              │
└──────────────┘    └──────────────┘    └──────────────┘
                            │
                            ▼
                    ┌──────────────┐
                    │ PostgreSQL   │
                    │ (Stats/Users)│
                    └──────────────┘
```

## Error Handling

```
User Input
    ↓
Validate Format ──→ Invalid? → Show error, retry
    ↓
Validate Range ──→ Out of bounds? → Show error, retry
    ↓
Check Available ──→ Taken? → Show error, retry
    ↓
Apply Move
    ↓
Update State
```

## Performance Characteristics

- **Latency**: < 10ms per move (local network)
- **Throughput**: Limited by SSH connections (~100-1000 concurrent)
- **Memory**: ~10MB base + ~1KB per game session
- **CPU**: Minimal (turn-based, no real-time processing)

## Security Model

### Current (Demo)
- No authentication
- No authorization
- No encryption beyond SSH
- No rate limiting
- Ephemeral rooms (no persistence)

### Production Recommendations
- SSH key authentication
- User management system
- Room access controls
- Input validation and sanitization
- Rate limiting per IP
- Audit logging
- Connection timeouts

## Scalability Considerations

### Current Limits
- Single server instance
- In-memory room storage
- No load balancing
- No session persistence

### Scaling Path
1. Add Redis for shared room state
2. Run multiple server instances
3. Add load balancer
4. Implement health checks
5. Add metrics/monitoring
6. Database for persistent data

## Testing Strategy

```
Unit Tests (board_test.go)
    ↓
Test Game Logic
    ├─ Board operations
    ├─ Win conditions
    ├─ Move validation
    └─ State management
```

### Test Coverage
- ✅ Board creation
- ✅ Move validation
- ✅ Win detection (all patterns)
- ✅ Draw detection
- ✅ Available moves
- ❌ AI logic (could add)
- ❌ Multiplayer sync (could add)
- ❌ Integration tests (could add)

## Build Process

```
Source Code (.go files)
        ↓
    go build
        ↓
Static Binary (6MB)
        ↓
   Docker Build
        ↓
Alpine Image (~15MB)
        ↓
   Container Registry
        ↓
  Deployment Target
```

## File Dependencies

```
main.go
  ├─→ pkg/server
  │     └─→ pkg/game
  │           ├─→ board.go
  │           ├─→ ai.go
  │           └─→ multiplayer.go
  │
  └─→ github.com/gliderlabs/ssh
      └─→ golang.org/x/crypto
```

---

This architecture balances:
- **Simplicity**: Easy to understand and modify
- **Performance**: Fast enough for the use case
- **Scalability**: Can be extended when needed
- **Maintainability**: Clean separation of concerns
