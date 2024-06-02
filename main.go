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
	board      [100]string
	players    []player
	directions [8]int
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

/* Determines if a square is valid */
func validSquare(i int) bool {
	return (i%10 >= 1 && i%10 <= 8)
}

func makeFreshBoard(b *board) *board {
	for i := 0; i < 100; i++ {
		b.board[i] = "*"
	}
	for i := 11; i < 89; i++ {
		if validSquare(i) {
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
	b := board{directions: [8]int{10, -10, 1, -1, 11, -11, 9, -9}}
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

/* Returns a string of the opponents piece */
func getOpp(piece string) string {
	var opp string
	if piece == "X" {
		opp = "O"
	} else {
		opp = "X"
	}
	return opp
}

/* Checks if a given move is valid in a direction */
func validMove(square, direction int, g game, p player) bool {
	if g.state.board[square] != "." {
		return false
	}
	opp := getOpp(p.piece)
	newSquare := square + direction
	if g.state.board[newSquare] != opp {
		return false
	}
	for validSquare(newSquare) {
		if g.state.board[newSquare] == "." {
			return false
		}
		if g.state.board[newSquare] == p.piece {
			return true
		}
		newSquare += direction
	}

	return false
}

func main() {

	p1, p2 := initializePlayer("Tayler", "O"), initializePlayer("Arbi", "X")
	gptr, bptr := initializeGame(*p1, *p2)
	printBoard(bptr.board)
	o, x := score(*gptr)
	fmt.Printf("SCORE:	White: %d Black: %d\n", o, x)
	fmt.Printf("PLAYERS: %s: %s, %s: %s\n", p1.name, p1.piece, p2.name, p2.piece)
	fmt.Printf("DIRECTIONS: %v\n", bptr.directions)
}
