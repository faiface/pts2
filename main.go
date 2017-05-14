package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type pieceName struct {
	player, number rune
}

type piece struct {
	pieceName
	position int
}

type gameState struct {
	length int
	pieces []piece
}

var (
	errNoSuchPiece = errors.New("no such piece")
	errKickOwn     = errors.New("can't kick own piece")
)

func newGame(sideLength int, pieces []piece) gameState {
	return gameState{
		length: sideLength * 4,
		pieces: pieces,
	}
}

func newRandomGame(sideLength int, players, numbers []rune) gameState {
	pieces := make([]piece, 0, len(players)*len(numbers))

	positions := rand.Perm(sideLength * 4)
	for _, player := range players {
		for _, number := range numbers {
			position := positions[0]
			positions = positions[1:]

			pieces = append(pieces, piece{
				pieceName{player, number},
				position,
			})
		}
	}

	return newGame(sideLength, pieces)
}

func (gs gameState) moved(name pieceName, amount int) (gameState, error) {
	toMove := -1
	for i, p := range gs.pieces {
		if p.pieceName == name {
			toMove = i
			break
		}
	}

	if toMove == -1 {
		return gameState{}, errNoSuchPiece
	}

	moved := gs.pieces[toMove]
	moved.position += amount
	moved.position %= gs.length

	result := gameState{length: gs.length}
	for i, p := range gs.pieces {
		if i == toMove {
			p = moved
		} else if p.position == moved.position && p.player == moved.player {
			return gameState{}, errKickOwn
		}
		if i == toMove || p.position != moved.position {
			result.pieces = append(result.pieces, p)
		}
	}

	return result, nil
}

func (gs gameState) numPieces(player rune) int {
	num := 0
	for _, p := range gs.pieces {
		if p.player == player {
			num++
		}
	}
	return num
}

func (gs gameState) numPlayers(players []rune) int {
	num := 0
	for _, player := range players {
		if gs.numPieces(player) > 0 {
			num++
		}
	}
	return num
}

func (gs gameState) String() string {
	side := gs.length / 4
	runeSide := 1 + (side+1)*2 + 1

	field := make([][]rune, runeSide)
	for i := range field {
		field[i] = make([]rune, runeSide)
		for j := range field[i] {
			field[i][j] = ' '
		}
	}

	for i := 0; i < runeSide; i++ {
		field[0][i] = '#'
		field[i][0] = '#'
		field[runeSide-1][i] = '#'
		field[i][runeSide-1] = '#'
		if i >= 3 && i < runeSide-3 {
			field[3][i] = '#'
			field[i][3] = '#'
			field[runeSide-4][i] = '#'
			field[i][runeSide-4] = '#'
		}
	}

	for _, p := range gs.pieces {
		var x, y int
		switch {
		case p.position < 1*side:
			pos := p.position - 0*side
			x, y = 1, 1+pos*2
		case p.position < 2*side:
			pos := p.position - 1*side
			x, y = 1+pos*2, runeSide-3
		case p.position < 3*side:
			pos := p.position - 2*side
			x, y = runeSide-3, runeSide-3-pos*2
		case p.position < 4*side:
			pos := p.position - 3*side
			x, y = runeSide-3-pos*2, 1
		default:
			panic("unreachable")
		}

		field[y][x] = p.pieceName.player
		field[y+1][x] = '\\'
		field[y][x+1] = '\\'
		field[y+1][x+1] = p.pieceName.number
	}

	var buf bytes.Buffer
	for i := range field {
		buf.WriteString(string(field[i]) + "\n")
	}
	return buf.String()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	players := []rune{'A', 'B', 'C', 'D'}
	numbers := []rune{'1', '2', '3'}

	gs := newRandomGame(10, players, numbers)

	for turn := 0; ; turn = (turn + 1) % len(players) {
		if gs.numPlayers(players) <= 1 {
			break
		}

		if gs.numPieces(players[turn]) == 0 {
			continue
		}

		fmt.Printf("%v", gs)
		fmt.Printf("%c's Turn!\n", players[turn])

		dice := rand.Intn(6) + 1
		fmt.Printf("Dice: %d\n", dice)

	again:
		var numberStr string
		fmt.Printf("Which piece to move? ")
		fmt.Scan(&numberStr)
		number := []rune(numberStr)[0]

		newGs, err := gs.moved(pieceName{players[turn], number}, dice)
		if err != nil {
			fmt.Printf("Invalid move: %s\n", err)
			goto again
		}

		gs = newGs
	}

	fmt.Print(gs)
	fmt.Println("Congratulations! You have won!")
	fmt.Println("==============================")
}
