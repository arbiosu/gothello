package logic

import (
	"fmt"

	"github.com/gothello/bots"
)

/* Represents a Player in a Game of Othello */
type Player struct {
	name  string
	piece string
}

/* Represents a Board in a Game of Othello */
type Board struct {
	Board      [100]string
	players    []Player
	directions [8]int
	turn       string
}

/* Captures Board state */
type history struct {
	moves []Board
}

/*
	 Represents a Game of Othello. O (white) and X (black) hold their respective
		scores
*/
type Game struct {
	state Board
	hist  history
	O     int
	X     int
}

func InitializePlayer(name, piece string) *Player {
	p := Player{name, piece}
	return &p
}

/* Determines if a square is valid */
func validSquare(i int) bool {
	return (i%10 >= 1 && i%10 <= 8)
}

func InitializeBoard(b *Board) *Board {
	for i := 0; i < 100; i++ {
		b.Board[i] = "*"
	}
	for i := 11; i < 89; i++ {
		if validSquare(i) {
			b.Board[i] = "."
		}
	}

	b.Board[44] = "O"
	b.Board[45] = "X"
	b.Board[54] = "X"
	b.Board[55] = "O"

	return b
}

func InitializeGame(p1, p2 Player) (*Game, *Board) {
	b := Board{directions: [8]int{10, -10, 1, -1, 11, -11, 9, -9}, turn: "X"}
	bptr := InitializeBoard(&b)
	b.players = append(b.players, p1, p2)
	history := history{}
	g := Game{b, history, 2, 2}

	return &g, bptr
}

/* Print the Board to the terminal */
func printBoard(b [100]string) {
	for i := 0; i < 100; i++ {
		if i%10 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("%v	", b[i])
	}
	fmt.Printf("\n")
}

/* Returns the current score of a Game. Updates the Game's score */
func score(g Game, bot bool) (int, int) {
	o, x := 0, 0
	for i := 11; i < 89; i++ {
		if g.state.Board[i] == "O" {
			o += 1
		}
		if g.state.Board[i] == "X" {
			x += 1
		}
	}
	if !bot {
		g.O = o
		g.X = x
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
func validMove(square, direction int, g Game, p Player) bool {
	if g.state.Board[square] != "." {
		return false
	}
	opp := getOpp(p.piece)
	newSquare := square + direction
	if g.state.Board[newSquare] != opp {
		return false
	}
	for validSquare(newSquare) {
		if g.state.Board[newSquare] == "." {
			return false
		}
		if g.state.Board[newSquare] == p.piece {
			return true
		}
		newSquare += direction
		// if newSquare is greater than last valid square. May need to revisit**
		if newSquare > 89 || newSquare < 11 {
			return true
		}
	}
	return false
}

/* Returns a slice of all valid directions for a given move */
func validDirections(square int, p Player, g Game) []int {
	valid := make([]int, 0, 8)
	for i := 0; i < len(g.state.directions); i++ {
		if validMove(square, g.state.directions[i], g, p) {
			valid = append(valid, g.state.directions[i])
		}
	}
	return valid
}

func availableMoves(p Player, g Game) map[int]bool {
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
Flips the pieces in all valid directions. Records state of the Board.
Updates Player turn.
*/
func flip(square int, p Player, g *Game) {
	opp := getOpp(p.piece)
	dirs := validDirections(square, p, *g)
	g.state.Board[square] = p.piece

	for i := 0; i < len(dirs); i++ {
		newSquare := square + dirs[i]
		for g.state.Board[newSquare] == opp {
			g.state.Board[newSquare] = p.piece
			newSquare += dirs[i]
		}
	}
	g.hist.moves = append(g.hist.moves, g.state)
	g.state.turn = opp
}

/* Used for bot play and calculation - does not affect actual state */
func flipStatic(square int, p Player, g Game) {
	opp := getOpp(p.piece)
	dirs := validDirections(square, p, g)
	g.state.Board[square] = p.piece

	for i := 0; i < len(dirs); i++ {
		newSquare := square + dirs[i]
		for g.state.Board[newSquare] == opp {
			g.state.Board[newSquare] = p.piece
			newSquare += dirs[i]
		}
	}
	g.state.turn = opp
}

/* Checks if a Game is over */
func gameOver(g Game) bool {
	o, x := availableMoves(g.state.players[0], g), availableMoves(g.state.players[1], g)
	if len(o) < 1 && len(x) < 1 {
		return true
	}
	return false
}

func gameType() int {
	var input int
	fmt.Printf("Enter 1 for to play against a human\n")
	fmt.Printf("Enter 2 to play against Randy, a bot that selects moves at random\n")
	fmt.Printf("Enter 3 to play against Max, a minimax bot...\n")
	fmt.Printf("Enter 4 for the rules of Othello\n")
	fmt.Printf("GAME TYPE: ")
	fmt.Scanln(&input)
	return input
}

func getPlayers() (*Player, *Player) {
	var name, name2 string
	fmt.Printf("Enter Player 1's name. They will be the \"O\" pieces: ")
	fmt.Scanln(&name)
	fmt.Printf("Enter Player 2's name. They will be the \"X\" pieces: ")
	fmt.Scanln(&name2)
	p1, p2 := InitializePlayer(name, "O"), InitializePlayer(name2, "X")
	return p1, p2
}

func getPlayer(turn string, g Game) Player {
	if turn == g.state.players[0].piece {
		return g.state.players[0]
	} else {
		return g.state.players[1]
	}
}

func result(gptr *Game) {
	o, x := score(*gptr, false)
	fmt.Println("------GAME OVER!------")
	printBoard(gptr.state.Board)
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

func gameStatus(gptr *Game) (Player, map[int]bool) {
	o, x := score(*gptr, true)
	curr := getPlayer(gptr.state.turn, *gptr)
	legal := availableMoves(curr, *gptr)

	fmt.Printf("CURRENT TURN: %v, %s\n", curr.name, gptr.state.turn)
	fmt.Printf("SCORE:	O: %d X: %d\n", o, x)
	fmt.Printf("LEGAL MOVES FOR %s: %v\n", gptr.state.turn, legal)
	return curr, legal
}

func humanMove(curr Player, legal map[int]bool, gptr *Game) {
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

func initRandy() (*Player, *Player) {
	var name string
	fmt.Printf("Initializing Randy...complete! Randy will play as \"O\"\n")
	fmt.Printf("Enter your name:")
	fmt.Scanln(&name)
	p1, p2 := InitializePlayer("Randy", "O"), InitializePlayer(name, "X")
	return p1, p2
}

func initMax() (*Player, *Player) {
	var name string
	fmt.Printf("Initializing Max...complete! Max will play as \"O\"\n")
	fmt.Printf("Enter your name:")
	fmt.Scanln(&name)
	p1, p2 := InitializePlayer("Max", "O"), InitializePlayer(name, "X")
	return p1, p2
}

func HumanGame() {
	p1, p2 := getPlayers()
	gptr, _ := InitializeGame(*p1, *p2)
	fmt.Printf("PLAYERS: %s: %s, %s: %s\n", p1.name, p1.piece, p2.name, p2.piece)
	for !gameOver(*gptr) {
		printBoard(gptr.state.Board)
		curr, legal := gameStatus(gptr)
		humanMove(curr, legal, gptr)
	}
	result(gptr)
}

func RandyGame() {
	p1, p2 := initRandy()
	gptr, _ := InitializeGame(*p1, *p2)
	fmt.Printf("PLAYERS: %s: %s, %s: %s\n", p1.name, p1.piece, p2.name, p2.piece)
	for !gameOver(*gptr) {
		printBoard(gptr.state.Board)
		curr, legal := gameStatus(gptr)
		if gptr.state.turn == "O" {
			fmt.Printf("Randy is thinking on a move...\n")
			move := bots.RandyMove(legal)
			fmt.Printf("Randy chose: %d", move)
			flip(move, curr, gptr)
		} else {
			humanMove(curr, legal, gptr)
		}
	}
	result(gptr)
}

func MaxGame() {
	p1, p2 := initMax()
	gptr, _ := InitializeGame(*p1, *p2)
	fmt.Printf("PLAYERS: %s: %s, %s: %s\n", p1.name, p1.piece, p2.name, p2.piece)
	for !gameOver(*gptr) {
		printBoard(gptr.state.Board)
		curr, legal := gameStatus(gptr)
		if gptr.state.turn == "O" {
			fmt.Printf("Max is thinking on a move...\n")
			move := bots.RandyMove(legal)
			fmt.Printf("Max chose: %d", move)
			flip(move, curr, gptr)
		} else {
			humanMove(curr, legal, gptr)
		}
	}
	result(gptr)
}

func displayHelpMsg() {
	fmt.Println("TODO: PRINT OTHELLO RULES")
}

func PlayGame() {
	gameType := gameType()
	switch {
	case gameType == 1:
		HumanGame()
	case gameType == 2:
		RandyGame()
	case gameType == 3:
		MaxGame()
	case gameType == 4:
		displayHelpMsg()
	default:
		fmt.Println("EXITING....")
	}
}
