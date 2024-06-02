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
	turn       string
}

/* Captures board state */
type history struct {
	moves []board
}

/*
	 Represents a game of Othello. O (white) and X (black) hold their respective
		scores
*/
type game struct {
	state board
	hist  history
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

func initializeBoard(b *board) *board {
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
	b := board{directions: [8]int{10, -10, 1, -1, 11, -11, 9, -9}, turn: "X"}
	bptr := initializeBoard(&b)
	b.players = append(b.players, p1, p2)
	history := history{}
	g := game{b, history, 2, 2}

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
		}
		if g.state.board[i] == "X" {
			x += 1
		}
	}
	g.O = o
	g.X = x

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

/* Returns a slice of all valid directions for a given move */
func validDirections(square int, p player, g game) []int {
	valid := make([]int, 0, 8)
	for i := 0; i < len(g.state.directions); i++ {
		if validMove(square, g.state.directions[i], g, p) {
			valid = append(valid, g.state.directions[i])
		}
	}
	return valid
}

func availableMoves(p player, g game) map[int]bool {
	// dirs := make([]int, 0)
	moves := make(map[int]bool)
	for i := 11; i < 89; i++ {
		if validSquare(i) {
			check := validDirections(i, p, g)
			if len(check) > 0 {
				moves[i] = true
			}
		}
	}
	return moves
}

/*
Flips the pieces in all valid directions. Records state of the board.
Updates player turn.
*/
func flip(square int, p player, g *game) {
	opp := getOpp(p.piece)
	dirs := validDirections(square, p, *g)
	g.state.board[square] = p.piece

	for i := 0; i < len(dirs); i++ {
		newSquare := square + dirs[i]
		for g.state.board[newSquare] == opp {
			g.state.board[newSquare] = p.piece
			newSquare += dirs[i]
		}
	}
	g.hist.moves = append(g.hist.moves, g.state)
	g.state.turn = opp
}

/* Checks if a game is over */
func gameOver(g game) bool {
	o, x := availableMoves(g.state.players[0], g), availableMoves(g.state.players[1], g)
	if len(o) < 1 && len(x) < 1 {
		return true
	}
	return false
}

func gameType() int {
	var input int
	fmt.Printf("Enter 1 for 1v1, 2 to play against a bot: ")
	fmt.Scanln(&input)
	return input
}

func getPlayers() (*player, *player) {
	var name, name2 string
	fmt.Printf("Enter player 1's name. They will be the \"O\" pieces: ")
	fmt.Scanln(&name)
	fmt.Printf("Enter player 1's name. They will be the \"X\" pieces: ")
	fmt.Scanln(&name2)
	p1, p2 := initializePlayer(name, "O"), initializePlayer(name2, "X")
	return p1, p2
}

func getPlayer(turn string, g game) player {
	if turn == g.state.players[0].piece {
		return g.state.players[0]
	} else {
		return g.state.players[1]
	}
}

func playGame() {
	// gameType := gameType()
	// if gameType == 1 {
	//	p1, p2 := getPlayers()
	//}
	p1, p2 := getPlayers()
	gptr, _ := initializeGame(*p1, *p2)
	fmt.Printf("PLAYERS: %s: %s, %s: %s\n", p1.name, p1.piece, p2.name, p2.piece)
	for !gameOver(*gptr) {
		printBoard(gptr.state.board)
		o, x := score(*gptr)
		curr := getPlayer(gptr.state.turn, *gptr)
		legal := availableMoves(curr, *gptr)
		fmt.Printf("CURRENT TURN: %v, %s\n", curr.name, gptr.state.turn)
		fmt.Printf("SCORE:	O: %d X: %d\n", o, x)
		fmt.Printf("LEGAL MOVES FOR %s: %v\n", gptr.state.turn, legal)
		var move int
		fmt.Printf("%v, enter your move. Choose any of the above legal moves: ", curr.name)
		fmt.Scanln(&move)
		elem, ok := legal[move]
		if !ok {
			fmt.Printf("%v INVALID MOVE!", elem)
		} else {
			flip(move, curr, gptr)
		}
	}
	o, x := score(*gptr)
	fmt.Println("------GAME OVER!------")
	printBoard(gptr.state.board)
	fmt.Printf("FINAL SCORE:	O: %d	X: %d\n", o, x)
	switch {
	case o > x:
		winner := getPlayer("O", *gptr)
		fmt.Printf("WINNER: %v. CONGRATULATIONS!\n", winner.name)
	case o < x:
		winner := getPlayer("X", *gptr)
		fmt.Printf("WINNER: %v. CONGRATULATIONS!\n", winner.name)
	default:
		fmt.Printf("RESULT: TIE GAME!")
	}
	fmt.Println("Thank you for playing Gothello!")
}

func main() {
	playGame()
}
