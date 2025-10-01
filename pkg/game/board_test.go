package game

import "testing"

func TestNewBoard(t *testing.T) {
	board := NewBoard()
	if board == nil {
		t.Fatal("NewBoard returned nil")
	}
	
	for i := 0; i < 9; i++ {
		if board.cells[i] != Empty {
			t.Errorf("Expected empty cell at position %d, got %v", i, board.cells[i])
		}
	}
}

func TestMakeMove(t *testing.T) {
	board := NewBoard()
	
	// Valid move
	if !board.MakeMove(0, X) {
		t.Error("Expected valid move to succeed")
	}
	
	// Invalid move (same position)
	if board.MakeMove(0, O) {
		t.Error("Expected invalid move to fail")
	}
	
	// Invalid position
	if board.MakeMove(10, X) {
		t.Error("Expected out of bounds move to fail")
	}
}

func TestCheckWinner(t *testing.T) {
	tests := []struct {
		name     string
		moves    []struct{ pos int; player Cell }
		expected Cell
		won      bool
	}{
		{
			name: "X wins row 1",
			moves: []struct{ pos int; player Cell }{
				{0, X}, {3, O}, {1, X}, {4, O}, {2, X},
			},
			expected: X,
			won:      true,
		},
		{
			name: "O wins column 2",
			moves: []struct{ pos int; player Cell }{
				{0, X}, {1, O}, {2, X}, {4, O}, {3, X}, {7, O},
			},
			expected: O,
			won:      true,
		},
		{
			name: "X wins diagonal",
			moves: []struct{ pos int; player Cell }{
				{0, X}, {1, O}, {4, X}, {2, O}, {8, X},
			},
			expected: X,
			won:      true,
		},
		{
			name: "Draw",
			moves: []struct{ pos int; player Cell }{
				{0, X}, {1, O}, {2, X},
				{3, O}, {4, X}, {5, X},
				{6, O}, {7, X}, {8, O},
			},
			expected: Empty,
			won:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := NewBoard()
			for _, move := range tt.moves {
				board.MakeMove(move.pos, move.player)
			}
			
			winner, won := board.CheckWinner()
			if won != tt.won {
				t.Errorf("Expected won=%v, got %v", tt.won, won)
			}
			if winner != tt.expected {
				t.Errorf("Expected winner=%v, got %v", tt.expected, winner)
			}
		})
	}
}

func TestIsFull(t *testing.T) {
	board := NewBoard()
	
	if board.IsFull() {
		t.Error("Empty board should not be full")
	}
	
	for i := 0; i < 9; i++ {
		board.MakeMove(i, X)
	}
	
	if !board.IsFull() {
		t.Error("Filled board should be full")
	}
}

func TestGetAvailableMoves(t *testing.T) {
	board := NewBoard()
	
	moves := board.GetAvailableMoves()
	if len(moves) != 9 {
		t.Errorf("Expected 9 available moves, got %d", len(moves))
	}
	
	board.MakeMove(0, X)
	board.MakeMove(4, O)
	
	moves = board.GetAvailableMoves()
	if len(moves) != 7 {
		t.Errorf("Expected 7 available moves, got %d", len(moves))
	}
}
