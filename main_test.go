package main

import "testing"

var (
	sideLength = 8

	pieces = []piece{
		{pieceName{'A', '1'}, 1},
		{pieceName{'B', '4'}, 2},
		{pieceName{'C', '8'}, 8},
		{pieceName{'A', '4'}, 9},
		{pieceName{'B', '2'}, 5},
		{pieceName{'B', '3'}, 4},
		{pieceName{'A', '9'}, 3},
	}
)

func TestNumPieces(t *testing.T) {
	table := []struct {
		player rune
		want   int
	}{
		{'A', 3},
		{'B', 3},
		{'C', 1},
	}

	gs := newGame(sideLength, pieces)

	for _, r := range table {
		if gs.numPieces(r.player) != r.want {
			t.Errorf("gs.numPieces(%q) = %v, want %v", r.player, gs.numPieces(r.player), r.want)
		}
	}
}

func TestNumPlayers(t *testing.T) {
	table := []struct {
		players []rune
		want    int
	}{
		{[]rune{'A', 'B', 'C', 'D'}, 3},
		{[]rune{'B', 'C', 'D'}, 2},
		{[]rune{'D', 'C', 'B', 'A'}, 3},
		{[]rune{'F', 'G', 'H', 'C'}, 1},
	}

	gs := newGame(sideLength, pieces)

	for _, r := range table {
		if gs.numPlayers(r.players) != r.want {
			t.Errorf("gs.numPlayers(%q) = %v, want %v", r.players, gs.numPlayers(r.players), r.want)
		}
	}
}

func TestSafeMove(t *testing.T) {
	table := []struct {
		name   pieceName
		amount int
	}{
		{pieceName{'C', '8'}, 4},
		{pieceName{'A', '1'}, 5},
		{pieceName{'B', '2'}, 3},
	}

	gs := newGame(sideLength, pieces)

	for _, r := range table {
		_, err := gs.moved(r.name, r.amount)
		if err != nil {
			t.Errorf("gs.moved(%q, %v) -> %v", r.name, r.amount, err)
		}
	}
}

func TestKickMove(t *testing.T) {
	table := []struct {
		name   pieceName
		amount int
		kicked pieceName
	}{
		{pieceName{'A', '1'}, 1, pieceName{'B', '4'}},
		{pieceName{'B', '4'}, 6, pieceName{'C', '8'}},
		{pieceName{'B', '2'}, 4, pieceName{'A', '4'}},
	}

	gs := newGame(sideLength, pieces)

	for _, r := range table {
		numBefore := gs.numPieces(r.kicked.player)
		newGs, err := gs.moved(r.name, r.amount)
		if err != nil {
			t.Errorf("gs.moved(%q, %v) -> %v", r.name, r.amount, err)
		}
		if newGs.numPieces(r.kicked.player) != numBefore-1 {
			t.Errorf("gs.numPieces(%q) = %v, want %v", r.kicked.player, newGs.numPieces(r.kicked.player), numBefore-1)
		}
	}
}

func TestKickOwn(t *testing.T) {
	table := []struct {
		name   pieceName
		amount int
	}{
		{pieceName{'A', '1'}, 2},
		{pieceName{'A', '9'}, 6},
		{pieceName{'B', '4'}, 3},
		{pieceName{'B', '3'}, 1},
	}

	gs := newGame(sideLength, pieces)

	for _, r := range table {
		_, err := gs.moved(r.name, r.amount)
		if err != errKickOwn {
			t.Errorf("gs.moved(%q, %v) -> %v, want %v", r.name, r.amount, err, errKickOwn)
		}
	}
}
