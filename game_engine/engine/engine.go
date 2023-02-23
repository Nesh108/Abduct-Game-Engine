package engine

import (
	"game_engine/types"
)

func NewGame(size uint) *types.Game {

	board := &types.Board{
		Board: make([][]*types.Unit, size),
	}
	for i := range board.Board {
		board.Board[i] = make([]*types.Unit, size)
	}

	player1 := CreatePlayer("A", types.Red, 6)
	player2 := CreatePlayer("B", types.Blue, 6)

	game := &types.Game{
		Board:         board,
		Players:       []*types.Player{player1, player2},
		CurrentPlayer: 0,
		Finished:      false,
		TurnNumber:    1,
	}

	// Position Player 1's units
	game.PositionUnit(types.WIZARD, player1, 4, 5)
	game.PositionUnit(types.DRAGON, player1, 3, 5)
	game.PositionUnit(types.DRAGON, player1, 3, 6)
	game.PositionUnit(types.DRAGON, player1, 4, 6)
	game.PositionUnit(types.KNIGHT, player1, 1, 6)
	game.PositionUnit(types.KNIGHT, player1, 2, 7)
	game.PositionUnit(types.KNIGHT, player1, 3, 8)

	// Position Player 2's units
	game.PositionUnit(types.WIZARD, player2, 5, 4)
	game.PositionUnit(types.DRAGON, player2, 5, 3)
	game.PositionUnit(types.DRAGON, player2, 6, 3)
	game.PositionUnit(types.DRAGON, player2, 6, 4)
	game.PositionUnit(types.KNIGHT, player2, 6, 1)
	game.PositionUnit(types.KNIGHT, player2, 7, 2)
	game.PositionUnit(types.KNIGHT, player2, 8, 3)

	return game
}

func CreatePlayer(name string, color types.Color, numHouses uint) *types.Player {
	return &types.Player{
		Name:      name,
		Color:     color,
		NumHouses: numHouses,
		Score:     0,
	}
}
