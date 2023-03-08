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
	fmt.Println(types.Red, "Abduct Game", types.Reset)
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

		isGameOver, errMove := game.HandleMove(move, args)
		if !isGameOver {
			if errMove != nil {
				game.Print()
				fmt.Printf(string(types.Red)+"Invalid move: %s\n"+types.Reset, errMove)
				continue
			} else {
				game.NextTurn()
				game.Print()
			}
		} else {
			game.Print()
			fmt.Printf(string(currentPlayer.Color)+"[Player %s]"+types.Reset+" WINS!", currentPlayer.Name)
			game.Finished = true
		}
	}
	// Implement the check mate mechanism
	// Implement the check mechanism
}

// Read a line of input from the terminal
func readLine(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
