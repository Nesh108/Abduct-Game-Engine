package main

import (
	"fmt"
	"game_engine/engine"
	"game_engine/types"
)

func main() {
	fmt.Println("====================================")
	fmt.Println(types.Red, "Ambush Game", types.Reset)
	fmt.Println("-----------------------------")

	game := engine.NewGame(10)

	game.Print()

	game.PositionHouse(0, 0)
	game.PositionHouse(0, 2)
	game.PositionHouse(0, 4)
	game.PositionHouse(0, 6)
	game.PositionHouse(0, 8)
	game.PositionHouse(2, 0)
	game.PositionHouse(2, 2)
	game.PositionHouse(2, 4)
	game.PositionHouse(2, 6)
	game.PositionHouse(2, 8)
	game.PositionHouse(4, 1)
	game.PositionHouse(4, 3)

	game.Print()

	// Implement Next Turn (automatic) - add ++ to turn number, Set game.CurrentPlayer and switch it after every PositionHouse/Implement the fighting mechanism
	// Implement the check mate mechanism
	// Implement the check mechanism
	// Implement the game loop asking each player to make their move (read orig pos and new pos)
	// Implement colors for players (so the board shows the different types)
	// Implement getting the available location for up/down/left/right movement
	// Add function to check if no available moves are available
	// Show how many points and houses each player has
	// Show the turn count
}
