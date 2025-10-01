package server

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"strings"

	"tictactoe-ssh/pkg/game"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func Start(port string) error {
	ssh.Handle(func(s ssh.Session) {
		handleSession(s)
	})

	// Generate ed25519 host key
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ed25519 key: %v", err)
	}

	// Convert to SSH format
	sshPrivateKey, err := gossh.NewSignerFromKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to create SSH signer: %v", err)
	}

	// Configure SSH server with ed25519 host key
	server := &ssh.Server{
		Addr:        ":" + port,
		HostSigners: []ssh.Signer{sshPrivateKey},
	}

	return server.ListenAndServe()
}

func handleSession(s ssh.Session) {
	defer s.Close()

	showWelcome(s)

	for {
		// Check if session is still active
		select {
		case <-s.Context().Done():
			return
		default:
		}

		fmt.Fprintf(s, "\nChoose game mode:\n")
		fmt.Fprintf(s, "1. Play against CPU\n")
		fmt.Fprintf(s, "2. Multiplayer (Create room)\n")
		fmt.Fprintf(s, "3. Multiplayer (Join room)\n")
		fmt.Fprintf(s, "4. Quit\n")

		// Read input with echo
		choice, err := readInputWithEcho(s, "\nEnter choice (1-4): ")
		if err != nil {
			// Check if session is still active
			select {
			case <-s.Context().Done():
				return
			default:
				return
			}
		}

		switch choice {
		case "1":
			playCPU(s)
		case "2":
			createMultiplayerRoom(s)
		case "3":
			joinMultiplayerRoom(s)
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

// readInputWithEcho reads input from the session and echoes it back
func readInputWithEcho(s ssh.Session, prompt string) (string, error) {
	fmt.Fprint(s, prompt)

	var input strings.Builder
	buf := make([]byte, 1)

	for {
		n, err := s.Read(buf)
		if err != nil || n == 0 {
			return "", err
		}

		char := buf[0]

		// Handle different input characters
		switch char {
		case '\n', '\r':
			// End of input
			fmt.Fprintf(s, "\n")
			return strings.TrimSpace(input.String()), nil
		case '\b', 127: // Backspace
			if input.Len() > 0 {
				// Remove last character from input buffer
				str := input.String()
				input.Reset()
				input.WriteString(str[:len(str)-1])
				// Send backspace sequence to terminal for visual feedback
				fmt.Fprintf(s, "\b \b")
			}
		case 3: // Ctrl+C
			return "", fmt.Errorf("interrupted")
		default:
			// Regular character
			if char >= 32 && char <= 126 { // Printable ASCII
				input.WriteByte(char)
				fmt.Fprint(s, string(char))
			}
		}
	}
}

func playCPU(s ssh.Session) {
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
		input, err := readInputWithEcho(s, "\nYour turn (X): ")
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
		cpuRow, cpuCol := cpuMove/3+1, cpuMove%3+1
		fmt.Fprintf(s, "\nCPU played: %d %d\n", cpuRow, cpuCol)
	}

	readInputWithEcho(s, "\nPress Enter to return to main menu...")
}

func createMultiplayerRoom(s ssh.Session) {
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

	// Wait for second player or session cancellation to avoid busy-waiting/leaks
	select {
	case <-room.OpponentJoined():
		fmt.Fprintf(s, "\nOpponent joined!\n")
	case <-s.Context().Done():
		// Creating player disconnected; clean up room and exit
		game.GlobalRoomManager.RemoveRoom(room.ID)
		return
	}

	playMultiplayer(s, room, player)
}

func joinMultiplayerRoom(s ssh.Session) {
	code, err := readInputWithEcho(s, "\nEnter room code: ")
	if err != nil {
		return
	}

	room, err := game.GlobalRoomManager.JoinRoom(code)
	if err != nil {
		fmt.Fprintf(s, "Error: %s\n", err.Error())
		readInputWithEcho(s, "Press Enter to continue...")
		return
	}

	player := &game.Player{
		Input:  make(chan string, 1),
		Output: make(chan string, 10),
	}

	symbol, err := room.AssignPlayer(player)
	if err != nil {
		fmt.Fprintf(s, "Error: %s\n", err.Error())
		readInputWithEcho(s, "Press Enter to continue...")
		return
	}

	fmt.Fprintf(s, "\n=== Joined Room %s ===\n", room.ID)
	fmt.Fprintf(s, "You are: %s\n", symbol.String())

	playMultiplayer(s, room, player)
}

func playMultiplayer(s ssh.Session, room *game.Room, player *game.Player) {
	defer game.GlobalRoomManager.RemoveRoom(room.ID)

	fmt.Fprintf(s, "\nEnter moves as: row col (e.g., '1 2' for row 1, column 2)\n")

	for {
		fmt.Fprintf(s, "%s\n", room.Board.String())
		fmt.Fprintf(s, "%s\n", room.GetStatus())

		// Get game state with proper synchronization
		current, gameOver, winner := room.GetGameState()

		if gameOver {
			switch winner {
			case player.Symbol:
				fmt.Fprintf(s, "\n🎉 Congratulations! You won!\n")
			case game.Empty:
				fmt.Fprintf(s, "\n🤝 It's a draw!\n")
			default:
				fmt.Fprintf(s, "\n😢 You lost! Better luck next time.\n")
			}
			break
		}

		if current == player.Symbol {
			input, err := readInputWithEcho(s, fmt.Sprintf("\nYour turn (%s): ", player.Symbol.String()))
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

			// Wait for opponent's move using channel-based synchronization
			select {
			case <-room.WaitForMove():
				// Move was made, continue to next iteration
			case <-s.Context().Done():
				// Session disconnected
				return
			}
		}
	}

	readInputWithEcho(s, "\nPress Enter to return to main menu...")
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
