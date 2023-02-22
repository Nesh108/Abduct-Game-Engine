package main

import (
	"fmt"
	"game_engine/engine"
	"game_engine/types"
)

func main() {
	fmt.Println("====================================")
	fmt.Println("\033[31mAbduct Game\033[0m")
	fmt.Println("-----------------------------")

	game := engine.NewGame(10)

	game.Print()

	game.PositionUnit(types.HOUSE, game.Players[0], 0, 0)
	game.PositionUnit(types.HOUSE, game.Players[1], 0, 2)
	game.PositionUnit(types.HOUSE, game.Players[0], 0, 4)
	game.PositionUnit(types.HOUSE, game.Players[1], 0, 6)
	game.PositionUnit(types.HOUSE, game.Players[0], 0, 8)
	game.PositionUnit(types.HOUSE, game.Players[1], 2, 0)
	game.PositionUnit(types.HOUSE, game.Players[0], 2, 2)
	game.PositionUnit(types.HOUSE, game.Players[1], 2, 4)
	game.PositionUnit(types.HOUSE, game.Players[0], 2, 6)
	game.PositionUnit(types.HOUSE, game.Players[1], 2, 8)
	game.PositionUnit(types.HOUSE, game.Players[0], 4, 1)
	game.PositionUnit(types.HOUSE, game.Players[1], 4, 3)

	game.Print()

	// Implement Next Turn (automatic) - add ++ to turn number, Set game.CurrentPlayer and switch it after every PositionUnit/etc.
	// Implement the fighting mechanism
	// Implement the check mate mechanism
	// Implement the check mechanism
	// Implement the game loop asking each player to make their move (read orig pos and new pos)
	// Implement colors for players (so the board shows the different types)
	// Implement getting the available location for up/down/left/right movement
	// Add function to check if no available moves are available
	// Show how many points and houses each player has
	// Show the turn count
}
