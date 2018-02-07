package tictactoe

import (
	"errors"
	"sync"

	"github.com/jinzhu/copier"
)

var ErrGameIsOver = errors.New("That game has finished.")
var ErrGameNotFound = errors.New("That game does not exist.")
var ErrInvalidMove = errors.New("That square on the board is occupied.")
var ErrInvalidPlayer = errors.New("Player is not in this game.")
var ErrNotPlayerTurn = errors.New("It is not this player's turn.")

type GameStatus uint

const (
	XMove GameStatus = iota
	OMove
	XWin
	OWin
	Draw
)

type Games struct {
	games  map[uint]*Game
	mutex  sync.RWMutex
	nextID uint
}

func NewGames() *Games {
	ret := &Games{}
	ret.games = make(map[uint]*Game)
	return ret
}

func (g *Games) GetGameIDs() []uint {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	ret := make([]uint, len(g.games))

	i := 0
	for k := range g.games {
		ret[i] = k
		i++
	}
	return ret
}

func (g *Games) NewGame(x uint, o uint) uint {
	game := &Game{}
	game.PlayerO = o
	game.PlayerX = x
	game.Board = [3][3]byte{
		{'-', '-', '-'},
		{'-', '-', '-'},
		{'-', '-', '-'},
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.games[g.nextID] = game
	ret := g.nextID
	g.nextID++

	return ret
}

func (g *Games) GetGame(id uint) (Game, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	game, ok := g.games[id]
	if !ok {
		return Game{}, ErrGameNotFound
	}
	ret := Game{}
	copier.Copy(&ret, &game)
	return ret, nil
}

func (g *Games) MakeMove(id, playerID uint, row, col byte) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	game, ok := g.games[id]
	if !ok {
		return ErrGameNotFound
	}
	return game.MakeMove(playerID, row, col)
}

type Game struct {
	Board     [3][3]byte
	PlayerO   uint
	PlayerX   uint
	Status    GameStatus
	moveCount uint
}

func (g *Game) MakeMove(playerID uint, row, col byte) error {
	// Make sure the game isn't over
	if g.Status == XWin || g.Status == OWin || g.Status == Draw {
		return ErrGameIsOver
	}

	// Player validation
	var player int
	if playerID == g.PlayerX {
		player = 0
	} else if playerID == g.PlayerO {
		player = 1
	} else {
		return ErrInvalidPlayer
	}

	// Confirm it is this player's turn
	if g.Status == XMove && player != 0 {
		return ErrNotPlayerTurn
	}
	if g.Status == OMove && player != 1 {
		return ErrNotPlayerTurn
	}

	// Make the move
	if g.Board[col][row] != '-' {
		return ErrInvalidMove
	}

	g.moveCount++
	switch player {
	case 0:
		g.Board[col][row] = 'x'
		g.Status = g.checkGameOver('x', row, col)
	case 1:
		g.Board[col][row] = 'o'
		g.Status = g.checkGameOver('o', row, col)
	}

	return nil
}

func (g *Game) checkGameOver(tok, row, col byte) GameStatus {
	// rows
	for i := 0; i < 3; i++ {
		if g.Board[col][i] != tok {
			break
		}
		if i == 2 {
			switch tok {
			case 'x':
				return XWin
			case 'o':
				return OWin
			}
		}
	}

	// cols
	for i := 0; i < 3; i++ {
		if g.Board[i][row] != tok {
			break
		}
		if i == 2 {
			switch tok {
			case 'x':
				return XWin
			case 'o':
				return OWin
			}
		}
	}

	// diag
	if row == col {
		for i := 0; i < 3; i++ {
			if g.Board[i][i] != tok {
				break
			}
			if i == 2 {
				switch tok {
				case 'x':
					return XWin
				case 'o':
					return OWin
				}
			}
		}
	}

	// anti-diag
	for i := 0; i < 3; i++ {
		if g.Board[i][2-i] != tok {
			break
		}
		if i == 2 {
			switch tok {
			case 'x':
				return XWin
			case 'o':
				return OWin
			}
		}
	}

	// Draw
	if g.moveCount == 8 {
		return Draw
	}

	if tok == 'x' {
		return OMove
	}
	return XMove
}
