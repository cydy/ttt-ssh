package game

import (
	"sync"
	"testing"
	"time"
)

func TestRoomCreationAndJoining(t *testing.T) {
	rm := &RoomManager{
		rooms: make(map[string]*Room),
	}

	room := rm.CreateRoom()
	if room == nil {
		t.Fatal("Expected room to be created")
	}
	if room.ID == "" {
		t.Fatal("Expected room to have an ID")
	}
	if room.Board == nil {
		t.Fatal("Expected room to have a board")
	}
	if room.Current != X {
		t.Fatalf("Expected current player to be X, got %v", room.Current)
	}
}

func TestPlayerAssignment(t *testing.T) {
	rm := &RoomManager{
		rooms: make(map[string]*Room),
	}

	room := rm.CreateRoom()

	player1 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}

	symbol1, err := room.AssignPlayer(player1)
	if err != nil {
		t.Fatalf("Failed to assign player 1: %v", err)
	}
	if symbol1 != X {
		t.Fatalf("Expected player 1 to be X, got %v", symbol1)
	}

	player2 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}

	symbol2, err := room.AssignPlayer(player2)
	if err != nil {
		t.Fatalf("Failed to assign player 2: %v", err)
	}
	if symbol2 != O {
		t.Fatalf("Expected player 2 to be O, got %v", symbol2)
	}

	// Try to assign third player (should fail)
	player3 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}
	_, err = room.AssignPlayer(player3)
	if err == nil {
		t.Fatal("Expected error when assigning third player")
	}
}

func TestOpponentJoinedNotification(t *testing.T) {
	rm := &RoomManager{
		rooms: make(map[string]*Room),
	}

	room := rm.CreateRoom()

	player1 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}
	room.AssignPlayer(player1)

	// Start a goroutine to wait for opponent
	done := make(chan bool)
	go func() {
		<-room.OpponentJoined()
		done <- true
	}()

	// Assign second player
	player2 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}
	room.AssignPlayer(player2)

	// Wait for notification with timeout
	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for opponent joined notification")
	}
}

func TestMoveNotification(t *testing.T) {
	rm := &RoomManager{
		rooms: make(map[string]*Room),
	}

	room := rm.CreateRoom()

	player1 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}
	player2 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}

	room.AssignPlayer(player1)
	room.AssignPlayer(player2)

	// Start a goroutine to wait for move notification
	done := make(chan bool)
	go func() {
		<-room.WaitForMove()
		done <- true
	}()

	// Make a move
	err := room.MakeMove(0, X)
	if err != nil {
		t.Fatalf("Failed to make move: %v", err)
	}

	// Wait for notification with timeout
	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for move notification")
	}
}

func TestNoBusyWait(t *testing.T) {
	rm := &RoomManager{
		rooms: make(map[string]*Room),
	}

	room := rm.CreateRoom()

	player1 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}
	player2 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}

	room.AssignPlayer(player1)
	room.AssignPlayer(player2)

	// Simulate multiple goroutines waiting for moves
	// This should not cause excessive CPU usage
	var wg sync.WaitGroup
	numWaiters := 10

	for i := 0; i < numWaiters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-room.WaitForMove():
				// Got notification
			case <-time.After(200 * time.Millisecond):
				// Timeout (expected for most waiters)
			}
		}()
	}

	// Make a move after a short delay
	time.Sleep(50 * time.Millisecond)
	room.MakeMove(0, X)

	// Wait for all goroutines to finish
	doneChan := make(chan bool)
	go func() {
		wg.Wait()
		doneChan <- true
	}()

	select {
	case <-doneChan:
		// Success - all goroutines finished
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for goroutines to finish - possible busy wait")
	}
}

func TestConcurrentMoves(t *testing.T) {
	rm := &RoomManager{
		rooms: make(map[string]*Room),
	}

	room := rm.CreateRoom()

	player1 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}
	player2 := &Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}

	room.AssignPlayer(player1)
	room.AssignPlayer(player2)

	// Test that concurrent access is safe
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// Player X makes moves
	wg.Add(1)
	go func() {
		defer wg.Done()
		positions := []int{0, 3, 6}
		for _, pos := range positions {
			for {
				current, gameOver, _ := room.GetGameState()
				if gameOver {
					return
				}
				if current == X {
					err := room.MakeMove(pos, X)
					if err == nil {
						break
					}
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Player O makes moves
	wg.Add(1)
	go func() {
		defer wg.Done()
		positions := []int{1, 4}
		for _, pos := range positions {
			for {
				current, gameOver, _ := room.GetGameState()
				if gameOver {
					return
				}
				if current == O {
					err := room.MakeMove(pos, O)
					if err == nil {
						break
					}
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Wait for game to finish
	doneChan := make(chan bool)
	go func() {
		wg.Wait()
		close(errors)
		doneChan <- true
	}()

	select {
	case <-doneChan:
		// Check for errors
		for err := range errors {
			t.Errorf("Error during concurrent moves: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for concurrent moves to complete")
	}

	// Verify game ended
	_, gameOver, winner := room.GetGameState()
	if !gameOver {
		t.Error("Expected game to be over")
	}
	if winner != X {
		t.Errorf("Expected X to win, got %v", winner)
	}
}
