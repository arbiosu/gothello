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

/*
	 Represents a game of Othello. O (white) and X (black) hold their respective
		scores
*/
type game struct {
	state board
	O     int
	X     int
}

func initializePlayer(name, piece string) *player {
	p := player{name, piece}
	return &p
}

/* Determines if a square is playable */
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
	g := game{b, 2, 2}

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

/* Returns the current score of a game. Updates the game's score */
func score(g game) (int, int) {
	o, x := 0, 0
	for i := 11; i < 89; i++ {
		if g.state.board[i] == "O" {
			o += 1
			g.O += 1
		}
		if g.state.board[i] == "X" {
			x += 1
			g.X += 1
		}
	}
	return o, x
}

func main() {

	p1, p2 := player{"Tayler", "O"}, player{"Arbi", "X"}
	g, bptr := initializeGame(p1, p2)
	printBoard(bptr.board)
	o, x := score(*g)
	fmt.Printf("SCORE:	White: %d Black: %d\n", o, x)
	fmt.Printf("PLAYERS: %s: %s, %s: %s\n", p1.name, p1.piece, p2.name, p2.piece)
}
