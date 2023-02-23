package main

import (
	"bufio"
	"fmt"
	"game_engine/engine"
	"game_engine/types"
	"os"
	"strings"
)

func main() {
	availableMoves := []types.Move{
		types.Move{ID: types.PUT_HOUSE, Command: "put-house", Arguments: "\t<position> (e.g. B1)"},
		types.Move{ID: types.MOVE_UNIT, Command: "move", Arguments: "\t<position> <direction> (up/down/left/right)"},
		types.Move{ID: types.PUT_HOUSE, Command: "ph", Arguments: "\t<position> (e.g. B1)"},
		types.Move{ID: types.MOVE_UNIT, Command: "m", Arguments: "\t<position> <direction> (up/down/left/right)"},
	}

	fmt.Println("====================================")
	fmt.Println(types.Red, "Ambush Game", types.Reset)
	fmt.Println("-----------------------------")

	game := engine.NewGame(10)
	game.Print()

	for !game.Finished {
		currentPlayer := game.Players[game.CurrentPlayer]
		fmt.Println()
		fmt.Printf(string(currentPlayer.Color)+"[Player %s]"+types.Reset+" > ", currentPlayer.Name)
		// Read input from player
		moveStr := readLine("Enter move: ")
		moveStr = strings.ToLower(moveStr)

		if moveStr == "help" {
			fmt.Println()
			fmt.Println("Available Moves: ")
			for _, move := range availableMoves {
				fmt.Println(move.Command, " ", move.Arguments)
			}
			fmt.Println()
			continue
		}

		// Parse move string into move struct
		move, args, errParseMove := engine.ParseMove(moveStr, availableMoves)
		if errParseMove != nil {
			game.Print()
			fmt.Printf(string(types.Red)+"Invalid move: %s\n"+types.Reset, errParseMove)
			continue
		}

		errMove := game.HandleMove(move, args)
		if errMove != nil {
			game.Print()
			fmt.Printf(string(types.Red)+"Invalid move: %s\n"+types.Reset, errMove)
			continue
		} else {
			game.NextTurn()
			game.Print()
		}
	}

	// game.PositionHouse(0, 0)
	// game.PositionHouse(0, 2)
	// game.PositionHouse(0, 4)
	// game.PositionHouse(0, 6)
	// game.PositionHouse(0, 8)
	// game.PositionHouse(2, 0)
	// game.PositionHouse(2, 2)
	// game.PositionHouse(2, 4)
	// game.PositionHouse(2, 6)
	// game.PositionHouse(2, 8)
	// game.PositionHouse(4, 1)
	// game.PositionHouse(4, 3)

	game.Print()

	// Implement the fighting mechanism
	// Implement the check mate mechanism
	// Implement the check mechanism
	// Implement getting the available location for up/down/left/right movement
	// Add function to check if no available moves are available
}

// Read a line of input from the terminal
func readLine(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
