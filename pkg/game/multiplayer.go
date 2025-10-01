package game

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Room struct {
	ID       string
	Board    *Board
	Player1  *Player
	Player2  *Player
	Current  Cell
	GameOver bool
	Winner   Cell
	mutex    sync.Mutex
	// opponentJoined is closed exactly once when Player2 successfully joins
	opponentJoined     chan struct{}
	opponentJoinedOnce sync.Once
	// moveNotifier is used to signal when a move has been made
	moveNotifier chan struct{}
}

type Player struct {
	Symbol Cell
	Input  chan string
	Output chan string
}

type RoomManager struct {
	rooms map[string]*Room
	mutex sync.RWMutex
}

var GlobalRoomManager = &RoomManager{
	rooms: make(map[string]*Room),
}

func (rm *RoomManager) CreateRoom() *Room {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	room := &Room{
		ID:             generateRoomCode(),
		Board:          NewBoard(),
		Current:        X,
		opponentJoined: make(chan struct{}),
		moveNotifier:   make(chan struct{}, 1),
	}
	rm.rooms[room.ID] = room
	return room
}

func (rm *RoomManager) JoinRoom(code string) (*Room, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	room, exists := rm.rooms[code]
	if !exists {
		return nil, fmt.Errorf("room not found")
	}
	if room.Player2 != nil {
		return nil, fmt.Errorf("room is full")
	}
	return room, nil
}

func (rm *RoomManager) RemoveRoom(code string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	delete(rm.rooms, code)
}

func generateRoomCode() string {
	// Generate a 6-character room code
	id := uuid.New().String()
	return id[:6]
}

func (r *Room) AssignPlayer(player *Player) (Cell, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.Player1 == nil {
		r.Player1 = player
		player.Symbol = X
		return X, nil
	} else if r.Player2 == nil {
		r.Player2 = player
		player.Symbol = O
		// Signal that the opponent has joined. Safe to call multiple times via Once.
		r.opponentJoinedOnce.Do(func() { close(r.opponentJoined) })
		return O, nil
	}
	return Empty, fmt.Errorf("room is full")
}

// OpponentJoined returns a channel that is closed when Player2 joins the room.
// Callers can select on this to be notified without busy-waiting.
func (r *Room) OpponentJoined() <-chan struct{} {
	return r.opponentJoined
}

// WaitForMove blocks until a move is made or the context is cancelled.
// Returns true if a move was made, false if context was cancelled.
func (r *Room) WaitForMove() <-chan struct{} {
	return r.moveNotifier
}

func (r *Room) MakeMove(pos int, player Cell) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.GameOver {
		return fmt.Errorf("game is over")
	}

	if r.Current != player {
		return fmt.Errorf("not your turn")
	}

	if !r.Board.MakeMove(pos, player) {
		return fmt.Errorf("invalid move")
	}

	// Check for winner
	if winner, won := r.Board.CheckWinner(); won {
		r.GameOver = true
		r.Winner = winner
	} else if r.Board.IsFull() {
		r.GameOver = true
		r.Winner = Empty
	} else {
		// Switch turns
		if r.Current == X {
			r.Current = O
		} else {
			r.Current = X
		}
	}

	// Notify waiting goroutines that a move has been made
	select {
	case r.moveNotifier <- struct{}{}:
	default:
		// Channel is full, notification already pending
	}

	return nil
}

func (r *Room) GetStatus() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	status := "Room: " + r.ID + "\n"
	status += "Players: "
	if r.Player1 != nil {
		status += "X "
	}
	if r.Player2 != nil {
		status += "O"
	}
	status += "\n"

	if !r.GameOver {
		status += fmt.Sprintf("Current turn: %s\n", r.Current.String())
	} else {
		if r.Winner != Empty {
			status += fmt.Sprintf("Winner: %s\n", r.Winner.String())
		} else {
			status += "Draw!\n"
		}
	}

	return status
}

// GetGameState returns the current game state (Current, GameOver, Winner) with proper locking
func (r *Room) GetGameState() (current Cell, gameOver bool, winner Cell) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.Current, r.GameOver, r.Winner
}
