package server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"tictactoe-ssh/pkg/game"

	"github.com/gliderlabs/ssh"
)

func Start(port string) error {
	ssh.Handle(func(s ssh.Session) {
		handleSession(s)
	})

	return ssh.ListenAndServe(":"+port, nil)
}

func handleSession(s ssh.Session) {
	defer s.Close()

	reader := bufio.NewReader(s)
	
	showWelcome(s)
	
	for {
		fmt.Fprintf(s, "\nChoose game mode:\n")
		fmt.Fprintf(s, "1. Play against CPU\n")
		fmt.Fprintf(s, "2. Multiplayer (Create room)\n")
		fmt.Fprintf(s, "3. Multiplayer (Join room)\n")
		fmt.Fprintf(s, "4. Quit\n")
		fmt.Fprintf(s, "\nEnter choice (1-4): ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			playCPU(s, reader)
		case "2":
			createMultiplayerRoom(s, reader)
		case "3":
			joinMultiplayerRoom(s, reader)
		case "4":
			fmt.Fprintf(s, "\nThanks for playing! Goodbye!\n")
			return
		default:
			fmt.Fprintf(s, "\nInvalid choice. Please try again.\n")
		}
	}
}

func showWelcome(w io.Writer) {
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "╔════════════════════════════════════╗\n")
	fmt.Fprintf(w, "║   Welcome to SSH Tic-Tac-Toe!     ║\n")
	fmt.Fprintf(w, "╚════════════════════════════════════╝\n")
}

func playCPU(s ssh.Session, reader *bufio.Reader) {
	board := game.NewBoard()
	ai := game.NewAI(game.O)
	playerSymbol := game.X

	fmt.Fprintf(s, "\n=== Playing against CPU ===\n")
	fmt.Fprintf(s, "You are: X\n")
	fmt.Fprintf(s, "CPU is: O\n")
	fmt.Fprintf(s, "\nEnter moves as: row col (e.g., '1 2' for row 1, column 2)\n")

	for {
		fmt.Fprintf(s, "%s\n", board.String())

		// Check for game over
		if winner, won := board.CheckWinner(); won {
			if winner == playerSymbol {
				fmt.Fprintf(s, "\n🎉 Congratulations! You won!\n")
			} else {
				fmt.Fprintf(s, "\n😢 CPU wins! Better luck next time.\n")
			}
			break
		}
		if board.IsFull() {
			fmt.Fprintf(s, "\n🤝 It's a draw!\n")
			break
		}

		// Player move
		fmt.Fprintf(s, "\nYour turn (X): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		pos := parseMove(input)
		if pos == -1 {
			fmt.Fprintf(s, "Invalid input. Use format: row col (e.g., '1 2')\n")
			continue
		}

		if !board.MakeMove(pos, playerSymbol) {
			fmt.Fprintf(s, "Invalid move. That position is already taken.\n")
			continue
		}

		// Check if player won
		if _, won := board.CheckWinner(); won {
			fmt.Fprintf(s, "%s\n", board.String())
			fmt.Fprintf(s, "\n🎉 Congratulations! You won!\n")
			break
		}
		if board.IsFull() {
			fmt.Fprintf(s, "%s\n", board.String())
			fmt.Fprintf(s, "\n🤝 It's a draw!\n")
			break
		}

		// CPU move
		cpuMove := ai.GetMove(board)
		board.MakeMove(cpuMove, game.O)
		row, col := cpuMove/3+1, cpuMove%3+1
		fmt.Fprintf(s, "\nCPU played: %d %d\n", row, col)
	}

	fmt.Fprintf(s, "\nPress Enter to return to main menu...")
	reader.ReadString('\n')
}

func createMultiplayerRoom(s ssh.Session, reader *bufio.Reader) {
	room := game.GlobalRoomManager.CreateRoom()
	
	player := &game.Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}

	symbol, _ := room.AssignPlayer(player)
	
	fmt.Fprintf(s, "\n=== Multiplayer Room Created ===\n")
	fmt.Fprintf(s, "Room Code: %s\n", room.ID)
	fmt.Fprintf(s, "You are: %s\n", symbol.String())
	fmt.Fprintf(s, "Waiting for opponent to join...\n")

	// Wait for second player
	go func() {
		for room.Player2 == nil {
			// Busy wait - in production, use a channel or condition variable
		}
		player.Output <- "Opponent joined!"
	}()

	// Wait for opponent message
	msg := <-player.Output
	fmt.Fprintf(s, "\n%s\n", msg)

	playMultiplayer(s, reader, room, player)
}

func joinMultiplayerRoom(s ssh.Session, reader *bufio.Reader) {
	fmt.Fprintf(s, "\nEnter room code: ")
	code, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	code = strings.TrimSpace(code)

	room, err := game.GlobalRoomManager.JoinRoom(code)
	if err != nil {
		fmt.Fprintf(s, "Error: %s\n", err.Error())
		fmt.Fprintf(s, "Press Enter to continue...")
		reader.ReadString('\n')
		return
	}

	player := &game.Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}

	symbol, err := room.AssignPlayer(player)
	if err != nil {
		fmt.Fprintf(s, "Error: %s\n", err.Error())
		fmt.Fprintf(s, "Press Enter to continue...")
		reader.ReadString('\n')
		return
	}

	fmt.Fprintf(s, "\n=== Joined Room %s ===\n", room.ID)
	fmt.Fprintf(s, "You are: %s\n", symbol.String())

	playMultiplayer(s, reader, room, player)
}

func playMultiplayer(s ssh.Session, reader *bufio.Reader, room *game.Room, player *game.Player) {
	defer game.GlobalRoomManager.RemoveRoom(room.ID)

	fmt.Fprintf(s, "\nEnter moves as: row col (e.g., '1 2' for row 1, column 2)\n")

	for {
		fmt.Fprintf(s, "%s\n", room.Board.String())
		fmt.Fprintf(s, "%s\n", room.GetStatus())

		if room.GameOver {
			if room.Winner == player.Symbol {
				fmt.Fprintf(s, "\n🎉 Congratulations! You won!\n")
			} else if room.Winner == game.Empty {
				fmt.Fprintf(s, "\n🤝 It's a draw!\n")
			} else {
				fmt.Fprintf(s, "\n😢 You lost! Better luck next time.\n")
			}
			break
		}

		if room.Current == player.Symbol {
			fmt.Fprintf(s, "\nYour turn (%s): ", player.Symbol.String())
			input, err := reader.ReadString('\n')
			if err != nil {
				return
			}

			pos := parseMove(input)
			if pos == -1 {
				fmt.Fprintf(s, "Invalid input. Use format: row col (e.g., '1 2')\n")
				continue
			}

			err = room.MakeMove(pos, player.Symbol)
			if err != nil {
				fmt.Fprintf(s, "Error: %s\n", err.Error())
				continue
			}
		} else {
			fmt.Fprintf(s, "\nWaiting for opponent's move...\n")
			
			// Poll for opponent's move
			currentBoard := room.Board.String()
			for {
				if room.Board.String() != currentBoard || room.GameOver {
					break
				}
			}
		}
	}

	fmt.Fprintf(s, "\nPress Enter to return to main menu...")
	reader.ReadString('\n')
}

func parseMove(input string) int {
	parts := strings.Fields(strings.TrimSpace(input))
	if len(parts) != 2 {
		return -1
	}

	row, err1 := strconv.Atoi(parts[0])
	col, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return -1
	}

	if row < 1 || row > 3 || col < 1 || col > 3 {
		return -1
	}

	return (row-1)*3 + (col - 1)
}
