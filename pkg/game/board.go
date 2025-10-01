package game

import (
	"fmt"
	"strings"
)

type Cell int

const (
	Empty Cell = iota
	X
	O
)

func (c Cell) String() string {
	switch c {
	case X:
		return "X"
	case O:
		return "O"
	default:
		return " "
	}
}

type Board struct {
	cells [9]Cell
}

func NewBoard() *Board {
	return &Board{}
}

func (b *Board) Reset() {
	b.cells = [9]Cell{}
}

func (b *Board) IsValidMove(pos int) bool {
	return pos >= 0 && pos < 9 && b.cells[pos] == Empty
}

func (b *Board) MakeMove(pos int, player Cell) bool {
	if !b.IsValidMove(pos) {
		return false
	}
	b.cells[pos] = player
	return true
}

func (b *Board) GetCell(pos int) Cell {
	if pos >= 0 && pos < 9 {
		return b.cells[pos]
	}
	return Empty
}

func (b *Board) IsFull() bool {
	for _, cell := range b.cells {
		if cell == Empty {
			return false
		}
	}
	return true
}

func (b *Board) CheckWinner() (Cell, bool) {
	// Check rows
	for i := 0; i < 3; i++ {
		if b.cells[i*3] != Empty &&
			b.cells[i*3] == b.cells[i*3+1] &&
			b.cells[i*3] == b.cells[i*3+2] {
			return b.cells[i*3], true
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if b.cells[i] != Empty &&
			b.cells[i] == b.cells[i+3] &&
			b.cells[i] == b.cells[i+6] {
			return b.cells[i], true
		}
	}

	// Check diagonals
	if b.cells[4] != Empty {
		if b.cells[0] == b.cells[4] && b.cells[4] == b.cells[8] {
			return b.cells[4], true
		}
		if b.cells[2] == b.cells[4] && b.cells[4] == b.cells[6] {
			return b.cells[4], true
		}
	}

	return Empty, false
}

func (b *Board) GetAvailableMoves() []int {
	moves := []int{}
	for i := 0; i < 9; i++ {
		if b.cells[i] == Empty {
			moves = append(moves, i)
		}
	}
	return moves
}

func (b *Board) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("     1   2   3\n")
	sb.WriteString("   ┌───┬───┬───┐\n")
	for row := 0; row < 3; row++ {
		sb.WriteString(fmt.Sprintf(" %d │", row+1))
		for col := 0; col < 3; col++ {
			pos := row*3 + col
			sb.WriteString(fmt.Sprintf(" %s │", b.cells[pos].String()))
		}
		sb.WriteString("\n")
		if row < 2 {
			sb.WriteString("   ├───┼───┼───┤\n")
		}
	}
	sb.WriteString("   └───┴───┴───┘\n")
	return sb.String()
}
