package main

import (
	"fmt"
)

/* Represents a player in a game of Othello */
type player struct {
	name  string
	piece string
}

/* Represents a board in a game of Othello */
type board struct {
	board   [100]string
	players []player
}

type game struct {
	state board
	score int
}

func initializePlayer(name, piece string) *player {
	p := player{name, piece}
	return &p
}

/* Gets the playable squares in an initial board */
func initialSquares(i int) bool {
	return (i%10 >= 1 && i%10 <= 8)
}

func makeFreshBoard(b *board) *board {
	for i := 0; i < 100; i++ {
		b.board[i] = "*"
	}
	for i := 11; i < 89; i++ {
		if initialSquares(i) {
			b.board[i] = "."
		}
	}

	b.board[44] = "O"
	b.board[45] = "X"
	b.board[54] = "X"
	b.board[55] = "O"

	return b
}

func initializeGame(p1, p2 player) (*game, *board) {
	b := board{}
	bptr := makeFreshBoard(&b)
	b.players = append(b.players, p1, p2)
	g := game{b, 0}

	return &g, bptr
}

/* Print the board to the terminal */
func printBoard(b [100]string) {
	for i := 0; i < 100; i++ {
		if i%10 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("%v	", b[i])
	}
	fmt.Printf("\n")
}

func main() {
	b := board{}
	x := makeFreshBoard(&b)
	printBoard(x.board)
}
