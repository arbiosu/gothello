package main

import "fmt"

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

func makeFreshBoard(board [100]string) [100]string {
	for i := 0; i < 100; i++ {
		board[i] = "*"
	}
	for i := 11; i < 89; i++ {
		if initialSquares(i) {
			board[i] = "."
		}
	}

	board[44] = "O"
	board[54] = "O"
	board[45] = "X"
	board[55] = "X"

	return board
}

func initializeGame(p1, p2 player) *game {
	b := board{}
	b.players = append(b.players, p1, p2)
	b.board = makeFreshBoard(b.board)
	g := game{b, 0}

	return &g
}

func main() {
	fmt.Println("Hello, Bastard!")
	b := [100]string{}
	x := makeFreshBoard(b)
	fmt.Printf("%v", x)
}
