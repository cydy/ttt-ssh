package game

import (
	"math/rand"
)

type AI struct {
	player Cell
}

func NewAI(player Cell) *AI {
	return &AI{player: player}
}

func (ai *AI) GetMove(board *Board) int {
	// Try to win
	if move := ai.findWinningMove(board, ai.player); move != -1 {
		return move
	}

	// Block opponent from winning
	opponent := O
	if ai.player == O {
		opponent = X
	}
	if move := ai.findWinningMove(board, opponent); move != -1 {
		return move
	}

	// Take center if available
	if board.IsValidMove(4) {
		return 4
	}

	// Take corners
	corners := []int{0, 2, 6, 8}
	availableCorners := []int{}
	for _, corner := range corners {
		if board.IsValidMove(corner) {
			availableCorners = append(availableCorners, corner)
		}
	}
	if len(availableCorners) > 0 {
		return availableCorners[rand.Intn(len(availableCorners))]
	}

	// Take any available move
	moves := board.GetAvailableMoves()
	if len(moves) > 0 {
		return moves[rand.Intn(len(moves))]
	}

	return -1
}

func (ai *AI) findWinningMove(board *Board, player Cell) int {
	for _, move := range board.GetAvailableMoves() {
		// Make a temporary move
		testBoard := *board
		testBoard.MakeMove(move, player)
		
		// Check if it's a winning move
		if winner, won := testBoard.CheckWinner(); won && winner == player {
			return move
		}
	}
	return -1
}
