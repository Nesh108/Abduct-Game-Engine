package engine

import (
	"errors"
	"fmt"
	"game_engine/types"
	"strings"
)

func NewGame(size uint) *types.Game {

	board := &types.Board{
		Board: make([][]*types.Unit, size),
	}
	for i := range board.Board {
		board.Board[i] = make([]*types.Unit, size)
	}

	player1 := CreatePlayer(0, "A", types.Red, 6)
	player2 := CreatePlayer(1, "B", types.Blue, 6)

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

func CreatePlayer(id uint, name string, color types.Color, numHouses uint) *types.Player {
	return &types.Player{
		ID:        id,
		Name:      name,
		Color:     color,
		NumHouses: numHouses,
		Score:     0,
	}
}

// Parse move string into a Move struct
func ParseMove(moveStr string, availableMoves []types.Move) (types.Move, []string, error) {
	// Split move string into command and args
	parts := strings.Split(moveStr, " ")
	command := parts[0]
	args := parts[1:]

	move, errMove := GetMove(availableMoves, command)
	// Check if command is valid
	if errMove != nil {
		return move, args, fmt.Errorf("invalid command: %s", command)
	}

	return move, args, nil
}

// Check if a slice of strings contains a given string
func GetMove(moves []types.Move, s string) (types.Move, error) {
	for _, m := range moves {
		if m.Command == s {
			return m, nil
		}
	}
	return types.Move{}, errors.New("move not found")
}
