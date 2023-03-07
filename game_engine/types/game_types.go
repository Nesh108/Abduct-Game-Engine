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

type Move struct {
	ID        MoveId
	Command   string
	Arguments string
}

type MoveId int

const (
	PUT_HOUSE MoveId = iota
	MOVE_UNIT
)

type Game struct {
	*Board                  // the 2D board of the game
	Players       []*Player // the two players
	CurrentPlayer uint      // the current player
	Finished      bool      // whether the game is finished
	TurnNumber    uint      // current turn id
}

type Player struct {
	ID        uint   // unique ID of the player
	Name      string // the name of the player
	Color     Color  // the color of the player (white or black)
	NumHouses uint   // the number of houses the player has left to place
	Score     uint   // current score
}

func (player *Player) Print() {
	fmt.Println(string(player.Color), "Player ", player.Name, ": ", player.Score, " | Houses: ", strconv.Itoa(int(player.NumHouses)), Reset)
}

type Unit struct {
	Name     UnitName // the name of the unit (wizard, dragon, knight or house)
	Color    Color    // the color of the unit (white or black)
	CanMove  bool     // whether the unit can move
	Position Position // the position of the unit on the board
	Owner    *Player  // who owns the piece
	Points   uint     // points given on kill
}

type Position struct {
	X uint // the x-coordinate of the position
	Y uint // the y-coordinate of the position
}

func StringToPosition(stringPos string, sizeBoard int) (Position, error) {
	pos := Position{}
	if len(stringPos) != 2 {
		return pos, fmt.Errorf("invalid position string: %s", stringPos)
	}
	stringPos = strings.ToUpper(stringPos)

	x, err := strconv.Atoi(string(stringPos[1]))
	if err != nil || x < 0 || x > 9 {
		return pos, fmt.Errorf("invalid y coordinate: %s", stringPos)
	}

	y := int(stringPos[0]) - int('A')
	if y < 0 || y >= sizeBoard {
		return pos, fmt.Errorf("invalid x coordinate: %s", stringPos)
	}

	pos.X = uint(x)
	pos.Y = uint(y)
	return pos, nil
}

type UnitName string

const (
	WIZARD UnitName = "Wizard"
	DRAGON UnitName = "Dragon"
	KNIGHT          = "Knight"
	HOUSE           = "House"
)

type Direction string

const (
	Up    Direction = "up"
	Down  Direction = "down"
	Left  Direction = "left"
	Right Direction = "right"
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

func (game *Game) NextTurn() {
	game.TurnNumber++

	// Swap player
	if game.CurrentPlayer == 0 {
		game.CurrentPlayer = 1
	} else {
		game.CurrentPlayer = 0
	}
}

func (game *Game) HandleMove(move Move, args []string) (bool, error) {
	if move.ID == PUT_HOUSE {
		pos, errPosition := StringToPosition(args[0], len(game.Board.Board[0]))
		if errPosition != nil {
			return false, errPosition
		}
		unit := game.Board.GetUnitAtPosition(pos)
		if unit != nil {
			return false, fmt.Errorf("cannot place house as unit already found at " + args[0])
		}
		return false, game.PositionHouse(pos.X, pos.Y)
	} else if move.ID == MOVE_UNIT {
		if len(args) != 2 {
			return false, fmt.Errorf("invalid command move (move <position> <direction>)")
		}

		pos, errPosition := StringToPosition(args[0], len(game.Board.Board[0]))
		if errPosition != nil {
			return false, errPosition
		}

		unit := game.Board.GetUnitAtPosition(pos)
		if unit == nil {
			return false, fmt.Errorf("no unit found at " + args[0])
		}

		dir, errDir := ParseDirection(args[1])
		if errDir != nil {
			return false, errDir
		}

		if unit.Owner.ID != game.Players[game.CurrentPlayer].ID {
			return false, fmt.Errorf("cannot move the opponent's piece")
		}
		return game.MoveUnit(unit, dir)
	}

	return false, fmt.Errorf("no move performed")
}

func ParseDirection(dirStr string) (Direction, error) {
	switch dirStr {
	case "up":
		return Up, nil
	case "down":
		return Down, nil
	case "left":
		return Left, nil
	case "right":
		return Right, nil
	default:
		return Up, fmt.Errorf("invalid direction (up, down, left, right)")
	}
}

func (game *Game) MoveUnit(unit *Unit, direction Direction) (bool, error) {
	switch unit.Name {
	case KNIGHT:
		return game.MoveKnight(unit, direction)
	case DRAGON:
		return game.MoveDragon(unit, direction)
	case WIZARD:
		return game.MoveWizard(unit, direction)
	case HOUSE:
		return false, fmt.Errorf("houses cannot be moved")
	}
	return false, fmt.Errorf("invalid unit selected: %s", unit.Name)
}

func (game *Game) MoveKnight(unit *Unit, direction Direction) (bool, error) {
	hasKilledWizard := false
	stepsMade := 0
	nextPosition := unit.Position.MoveByOneTowardsDirection(direction)
	for game.Board.IsValid(nextPosition) {
		otherUnit := game.Board.GetUnitAtPosition(nextPosition)
		if otherUnit != nil {
			// Unit from someone else
			if otherUnit.Owner.ID != game.Players[game.CurrentPlayer].ID {
				if otherUnit.Name != HOUSE {
					game.KillUnit(otherUnit)
					game.ChangeUnitPosition(unit, nextPosition)
					stepsMade++
					// Game over if wizard is killed
					if otherUnit.Name == WIZARD {
						hasKilledWizard = true
						break
					}
					break
				} else {
					// Other team house - Stopping
					break
				}
			} else {
				// Unit from own team
				break
			}
		}

		// Nothing on the way, keeps going
		stepsMade++
		game.ChangeUnitPosition(unit, nextPosition)

		nextPosition = unit.Position.MoveByOneTowardsDirection(direction)
	}

	if stepsMade > 0 {
		return hasKilledWizard, nil
	} else {
		return hasKilledWizard, fmt.Errorf("Unit cannot be moved towards that direction")
	}
}

func (game *Game) ChangeUnitPosition(unit *Unit, newPosition Position) {
	game.Board.Board[unit.Position.X][unit.Position.Y] = nil
	game.Board.Board[newPosition.X][newPosition.Y] = unit
	unit.Position = newPosition
}

func (game *Game) MoveDragon(unit *Unit, direction Direction) (bool, error) {
	hasKilledWizard := false
	hasAlreadyKilledDragon := false
	stepsMade := 0
	nextPosition := unit.Position.MoveByOneTowardsDirection(direction)
	for game.Board.IsValid(nextPosition) {
		otherUnit := game.Board.GetUnitAtPosition(nextPosition)
		if otherUnit != nil {
			// Unit from someone else
			if otherUnit.Owner.ID != game.Players[game.CurrentPlayer].ID {
				if otherUnit.Name != HOUSE {
					if otherUnit.Name != DRAGON {
						// Any other unit is ok
						game.KillUnit(otherUnit)
					} else if otherUnit.Name == DRAGON && !hasAlreadyKilledDragon {
						// Can only kill 1 dragon per swoop
						hasAlreadyKilledDragon = true
						game.KillUnit(otherUnit)
					}
					// Game over if wizard is killed
					if otherUnit.Name == WIZARD {
						hasKilledWizard = true
						break
					}
					// Keep going
				} else {
					// Other team house - Stopping
					break
				}
			} else {
				// Unit from own team
				break
			}
		}

		// Nothing on the way, keeps going
		stepsMade++
		game.ChangeUnitPosition(unit, nextPosition)

		nextPosition = unit.Position.MoveByOneTowardsDirection(direction)
	}

	if stepsMade > 0 {
		return hasKilledWizard, nil
	} else {
		return hasKilledWizard, fmt.Errorf("Unit cannot be moved towards that direction")
	}
}

func (game *Game) MoveWizard(unit *Unit, direction Direction) (bool, error) {
	hasKilledWizard := false
	stepsMade := 0
	nextPosition := unit.Position.MoveByOneTowardsDirection(direction)
	for game.Board.IsValid(nextPosition) {
		otherUnit := game.Board.GetUnitAtPosition(nextPosition)
		if otherUnit != nil {
			// Unit from someone else
			if otherUnit.Owner.ID != game.Players[game.CurrentPlayer].ID {
				if otherUnit.Name != HOUSE {
					// Any unit is ok
					game.KillUnit(otherUnit)

					// Game over if wizard is killed
					if otherUnit.Name == WIZARD {
						hasKilledWizard = true
						break
					}
					// else keep going
				} else {
					// Other team house - Stopping
					break
				}
			} else {
				// Unit from own team
				break
			}
		}

		// Nothing on the way, keeps going
		stepsMade++
		game.ChangeUnitPosition(unit, nextPosition)

		nextPosition = unit.Position.MoveByOneTowardsDirection(direction)
	}

	if stepsMade > 0 {
		return hasKilledWizard, nil
	} else {
		return hasKilledWizard, fmt.Errorf("Unit cannot be moved towards that direction")
	}
}

func (p Position) MoveByOneTowardsDirection(direction Direction) Position {
	switch direction {
	case Up:
		return Position{X: p.X - 1, Y: p.Y}
	case Left:
		return Position{X: p.X, Y: p.Y - 1}
	case Right:
		return Position{X: p.X, Y: p.Y + 1}
	case Down:
		return Position{X: p.X + 1, Y: p.Y}
	default:
		return p
	}
}

func (game *Game) KillUnit(unit *Unit) {
	game.Players[game.CurrentPlayer].Score += unit.Points

	game.Board.Board[unit.Position.X][unit.Position.Y] = nil
}

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

func (board *Board) IsValidPosition(position Position, unitName UnitName) error {
	if !board.IsValid(position) {
		return fmt.Errorf("position out of the board")
	}

	// Check if there is already a unit at the specified position
	if board.GetUnitAtPosition(position) != nil {
		return fmt.Errorf("there is already a unit at the specified position")
	}

	// Check if there is a house next to the specified position
	if unitName == HOUSE && board.IsAdjacentToHouse(position) {
		return fmt.Errorf("cannot place unit next to a house")
	}

	return nil
}

func (game *Game) PositionUnit(unitName UnitName, player *Player, x, y, points uint) error {
	if unitName == HOUSE {
		return errors.New("cannot place houses with this")
	}

	pos := Position{X: x, Y: y}

	// Check if the position is valid
	posErr := game.Board.IsValidPosition(pos, unitName)
	if posErr != nil {
		return posErr
	}

	unit := &Unit{
		Name:     unitName,
		Color:    player.Color,
		CanMove:  true,
		Position: pos,
		Owner:    player,
		Points:   points,
	}

	game.Board.Board[x][y] = unit

	return nil
}

func (game *Game) PositionHouse(x, y uint) error {
	player := game.Players[game.CurrentPlayer]
	if player.NumHouses <= 0 {
		return errors.New("no more houses available for placement")
	}

	pos := Position{X: x, Y: y}

	posErr := game.Board.IsValidPosition(pos, HOUSE)
	// Check if the position is valid
	if posErr != nil {
		return posErr
	}

	player.NumHouses--

	unit := &Unit{
		Name:     HOUSE,
		Color:    player.Color,
		CanMove:  false,
		Position: pos,
		Owner:    player,
	}

	game.Board.Board[x][y] = unit

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
	fmt.Println()

	// Print column labels
	fmt.Print(" " + VERTICAL_TILE)
	for col := 0; col < len(board.Board[0]); col++ {
		fmt.Printf(" %c ", 'A'+col)
		fmt.Print(VERTICAL_TILE)
	}
	fmt.Println()
	// Print top border
	fmt.Println(strings.Repeat(BORDER_TILE, len(board.Board[0])*2+1))

	for row, unitRow := range board.Board {
		// Print left border
		fmt.Print(strconv.Itoa(row))
		fmt.Print(VERTICAL_TILE)

		for _, unit := range unitRow {
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
		fmt.Println(strings.Repeat(BORDER_TILE, len(board.Board[0])*2+1))
	}
}
