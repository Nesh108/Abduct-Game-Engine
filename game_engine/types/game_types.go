package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Board struct {
	Board [][]*Unit
}

type Game struct {
	*Board                  // the 2D board of the game
	Players       []*Player // the two players
	CurrentPlayer *Player   // the current player
	Finished      bool      // whether the game is finished
	TurnNumber    uint      // current turn id
}

type Player struct {
	Name      string // the name of the player
	Color     Color  // the color of the player (white or black)
	NumHouses uint   // the number of houses the player has left to place
	Score     uint   // current score
}

func (player *Player) Print() {
	fmt.Println("Player ", player.Name, ": ", player.Score, " | Houses: ", strconv.Itoa(int(player.NumHouses)))
}

type Unit struct {
	Name     UnitName // the name of the unit (wizard, dragon, knight or house)
	Color    Color    // the color of the unit (white or black)
	CanMove  bool     // whether the unit can move
	Position Position // the position of the unit on the board
	Owner    *Player  // who owns the piece
}

type Position struct {
	X uint // the x-coordinate of the position
	Y uint // the y-coordinate of the position
}

type UnitName int

const (
	WIZARD UnitName = iota
	DRAGON
	KNIGHT
	HOUSE
)

type Color string

const (
	Reset  string = "\033[0m"
	Red    Color  = "\033[31m"
	Green  Color  = "\033[32m"
	Yellow Color  = "\033[33m"
	Blue   Color  = "\033[34m"
	Purple Color  = "\033[35m"
	Cyan   Color  = "\033[36m"
	Gray   Color  = "\033[37m"
	Bold   string = "\033[1m"
)

func (board *Board) Adjacent(pos Position) []Position {
	adjacent := []Position{
		{pos.X - 1, pos.Y - 1},
		{pos.X - 1, pos.Y},
		{pos.X - 1, pos.Y + 1},
		{pos.X, pos.Y - 1},
		{pos.X, pos.Y + 1},
		{pos.X + 1, pos.Y - 1},
		{pos.X + 1, pos.Y},
		{pos.X + 1, pos.Y + 1},
	}

	// Remove any positions that are out of bounds
	validPositions := []Position{}
	for _, p := range adjacent {
		if board.IsValid(p) {
			validPositions = append(validPositions, p)
		}
	}

	return validPositions
}

func (board *Board) IsValid(position Position) bool {
	// Check that the row and column are within the bounds of the board
	return position.X < uint(len(board.Board)) && position.Y < uint(len(board.Board))
}

func (board *Board) GetUnitAtPosition(position Position) *Unit {
	return board.Board[position.X][position.Y]
}

func (board *Board) IsAdjacentToHouse(position Position) bool {
	// check all adjacent positions to see if any are houses
	for _, adjPos := range board.Adjacent(position) {
		unit := board.GetUnitAtPosition(adjPos)
		if unit != nil && unit.Name == HOUSE {
			return true
		}
	}
	// no adjacent house found
	return false
}

func (board *Board) IsValidPosition(position Position, unitName UnitName) bool {
	if !board.IsValid(position) {
		return false
	}

	// Check if there is already a unit at the specified position
	if board.GetUnitAtPosition(position) != nil {
		fmt.Println("there is already a unit at the specified position")
		return false
	}

	// Check if there is a house next to the specified position
	if unitName == HOUSE && board.IsAdjacentToHouse(position) {
		fmt.Println("cannot place unit next to a house")
		return false
	}

	return true
}

func (board *Board) PositionUnit(unitName UnitName, player *Player, x, y uint, canMove ...bool) error {
	if unitName == HOUSE && player.NumHouses <= 0 {
		return errors.New("no more houses available for placement")
	}

	pos := Position{X: x, Y: y}

	// Check if the position is valid
	if !board.IsValidPosition(pos, unitName) {
		return errors.New("position is invalid")
	}

	if unitName == HOUSE {
		player.NumHouses--
	}

	unit := &Unit{
		Name:     unitName,
		Color:    player.Color,
		CanMove:  true,
		Position: pos,
		Owner:    player,
	}
	if len(canMove) > 0 {
		unit.CanMove = canMove[0]
	}

	board.Board[x][y] = unit

	return nil
}

const (
	WIZARD_TILE   = "W"
	DRAGON_TILE   = "D"
	KNIGHT_TILE   = "K"
	HOUSE_TILE    = "H"
	EMPTY_TILE    = " "
	BORDER_TILE   = "--"
	VERTICAL_TILE = "|"
)

func (game *Game) Print() {
	fmt.Println("Turn: ", game.TurnNumber)
	game.Players[0].Print()
	game.Players[1].Print()
	game.Board.Print()
}

func (board *Board) Print() {
	// Print top border
	fmt.Println(strings.Repeat(BORDER_TILE, 10*2+1))

	for _, row := range board.Board {
		// Print left border
		fmt.Print(VERTICAL_TILE)

		for _, unit := range row {
			var tile string

			if unit == nil {
				tile = EMPTY_TILE
			} else {
				switch unit.Name {
				case WIZARD:
					tile = Bold + (string(unit.Owner.Color) + WIZARD_TILE) + Reset
				case DRAGON:
					tile = Bold + (string(unit.Owner.Color) + DRAGON_TILE) + Reset
				case KNIGHT:
					tile = Bold + (string(unit.Owner.Color) + KNIGHT_TILE) + Reset
				case HOUSE:
					tile = Bold + (string(unit.Owner.Color) + HOUSE_TILE) + Reset
				}
			}

			fmt.Printf(" %s ", tile)
			fmt.Print(VERTICAL_TILE)
		}

		fmt.Println()

		// Print horizontal border
		fmt.Println(strings.Repeat(BORDER_TILE, 10*2+1))
	}
}
