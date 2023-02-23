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
	fmt.Println("====================================")
	fmt.Println(types.Red, "Ambush Game", types.Reset)
	fmt.Println("-----------------------------")

	game := engine.NewGame(10)

	availableMoves := []types.Move{
		types.Move{ID: types.PUT_HOUSE, Command: "put-house"},
		types.Move{ID: types.MOVE_UNIT, Command: "move"},
	}

	for !game.Finished {
		game.Print()

		fmt.Printf("Player %s, enter your move: ", game.Players[game.CurrentPlayer].Name)
		// Read input from player
		moveStr := readLine("Enter move: ")

		// Parse move string into move struct
		move, args, errParseMove := engine.ParseMove(moveStr, availableMoves)
		if errParseMove != nil {
			fmt.Printf("Invalid move: %s\n", errParseMove)
			continue
		}

		errMove := game.HandleMove(move, args)
		if errMove != nil {
			fmt.Printf("Invalid move: %s\n", errMove)
			continue
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
