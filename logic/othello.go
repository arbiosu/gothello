package logic

import (
	"fmt"
)

/* Represents a Player in a Game of Othello */
type Player struct {
	Name  string
	Piece string
}

/* Represents a Board in a Game of Othello */
type Board struct {
	Board      [100]string
	players    []Player
	directions [8]int
	Turn       string
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
	State Board
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
	b := Board{directions: [8]int{10, -10, 1, -1, 11, -11, 9, -9}, Turn: "X"}
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
		if g.State.Board[i] == "O" {
			o += 1
		}
		if g.State.Board[i] == "X" {
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
	if g.State.Board[square] != "." {
		return false
	}
	opp := getOpp(p.Piece)
	newSquare := square + direction
	if g.State.Board[newSquare] != opp {
		return false
	}
	for validSquare(newSquare) {
		if g.State.Board[newSquare] == "." {
			return false
		}
		if g.State.Board[newSquare] == p.Piece {
			return true
		}
		newSquare += direction
		// TODO: may be able to get away with removing this?
		if g.State.Board[newSquare] == "*" {
			return false
		}
	}
	return false
}

/* Returns a slice of all valid directions for a given move */
func validDirections(square int, p Player, g Game) []int {
	valid := make([]int, 0, 8)
	for i := 0; i < len(g.State.directions); i++ {
		if validMove(square, g.State.directions[i], g, p) {
			valid = append(valid, g.State.directions[i])
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
Updates Player Turn.
*/
func Flip(square int, p Player, g *Game) {
	opp := getOpp(p.Piece)
	if !validSquare(square) {
		if square == -1 {
			g.State.Turn = opp
		}
		return
	}
	dirs := validDirections(square, p, *g)
	g.State.Board[square] = p.Piece

	for i := 0; i < len(dirs); i++ {
		newSquare := square + dirs[i]
		for g.State.Board[newSquare] == opp {
			g.State.Board[newSquare] = p.Piece
			newSquare += dirs[i]
		}
	}
	g.hist.moves = append(g.hist.moves, g.State)
	g.State.Turn = opp
}

// Used for bot play and calculation - does not affect actual state
func flipStatic(square int, p Player, g Game) {
	opp := getOpp(p.Piece)
	dirs := validDirections(square, p, g)
	g.State.Board[square] = p.Piece

	for i := 0; i < len(dirs); i++ {
		newSquare := square + dirs[i]
		for g.State.Board[newSquare] == opp {
			g.State.Board[newSquare] = p.Piece
			newSquare += dirs[i]
		}
	}
}

/* Checks if a Game is over */
func GameOver(g Game) bool {
	o, x := availableMoves(g.State.players[0], g), availableMoves(g.State.players[1], g)
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
	if turn == g.State.players[0].Piece {
		return g.State.players[0]
	} else {
		return g.State.players[1]
	}
}

func result(gptr *Game) {
	o, x := score(*gptr, false)
	fmt.Println("------GAME OVER!------")
	printBoard(gptr.State.Board)
	fmt.Printf("FINAL SCORE:	O: %d	X: %d\n", o, x)
	switch {
	case o > x:
		winner := getPlayer("O", *gptr)
		fmt.Printf("WINNER: %v. CONGRATULATIONS!\n", winner.Name)
	case o < x:
		winner := getPlayer("X", *gptr)
		fmt.Printf("WINNER: %v. CONGRATULATIONS!\n", winner.Name)
	default:
		fmt.Printf("RESULT: TIE GAME!")
	}
	fmt.Println("Thank you for playing Gothello!")
}

func gameStatus(gptr *Game) (Player, map[int]bool) {
	o, x := score(*gptr, true)
	curr := getPlayer(gptr.State.Turn, *gptr)
	legal := availableMoves(curr, *gptr)

	fmt.Printf("CURRENT TURN: %v, %s\n", curr.Name, gptr.State.Turn)
	fmt.Printf("SCORE:	O: %d X: %d\n", o, x)
	fmt.Printf("LEGAL MOVES FOR %s: %v\n", gptr.State.Turn, legal)
	return curr, legal
}

func OnlineGameStatus(g *Game) (Player, map[int]bool, int, int) {
	o, x := score(*g, true)
	curr := getPlayer(g.State.Turn, *g)
	legal := availableMoves(curr, *g)
	return curr, legal, o, x
}

func humanMove(curr Player, legal map[int]bool, gptr *Game) {
	var move int
	fmt.Printf("%v, enter your move. Choose any of the above legal moves: ", curr.Name)
	fmt.Scanln(&move)
	elem, ok := legal[move]
	if !ok {
		fmt.Printf("%v INVALID MOVE!", elem)
	} else {
		Flip(move, curr, gptr)
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
	fmt.Printf("PLAYERS: %s: %s, %s: %s\n", p1.Name, p1.Piece, p2.Name, p2.Piece)
	for !GameOver(*gptr) {
		printBoard(gptr.State.Board)
		curr, legal := gameStatus(gptr)
		humanMove(curr, legal, gptr)
	}
	result(gptr)
}

func RandyGame() {
	p1, p2 := initRandy()
	gptr, _ := InitializeGame(*p1, *p2)
	fmt.Printf("PLAYERS: %s: %s, %s: %s\n", p1.Name, p1.Piece, p2.Name, p2.Piece)
	for !GameOver(*gptr) {
		printBoard(gptr.State.Board)
		curr, legal := gameStatus(gptr)
		if gptr.State.Turn == "O" {
			fmt.Printf("Randy is thinking on a move...\n")
			move := RandyMove(legal)
			fmt.Printf("Randy chose: %d", move)
			Flip(move, curr, gptr)
		} else {
			humanMove(curr, legal, gptr)
		}
	}
	result(gptr)
}

func MaxGame() {
	p1, p2 := initMax()
	gptr, _ := InitializeGame(*p1, *p2)
	fmt.Printf("PLAYERS: %s: %s, %s: %s\n", p1.Name, p1.Piece, p2.Name, p2.Piece)
	for !GameOver(*gptr) {
		printBoard(gptr.State.Board)
		curr, legal := gameStatus(gptr)
		if gptr.State.Turn == "O" {
			fmt.Printf("Max is thinking on a move...\n")
			move := MaxMove(gptr, curr, 3)
			fmt.Printf("Max chose: %d", move)
			Flip(move, curr, gptr)
		} else {
			humanMove(curr, legal, gptr)
		}
	}
	result(gptr)
}

func displayHelpMsg() {
	fmt.Println("TODO: PRINT OTHELLO RULES")
}

// Play Othello in the terminal!
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

func heuristic(g *Game, p Player) int {
	o, x := 0, 0
	for i := 11; i < 89; i++ {
		if g.State.Board[i] == "O" {
			mappedIndex := mapToStaticWeights(i)
			o += StaticWeights[mappedIndex]
		}
		if g.State.Board[i] == "X" {
			mappedIndex := mapToStaticWeights(i)
			x += StaticWeights[mappedIndex]
		}
	}
	if p.Piece == "O" {
		return o - x
	}
	return x - o
}

func minimax(g Game, depth int, isMax bool, p Player) int {
	if GameOver(g) || depth == 0 {
		return heuristic(&g, p)
	}
	legalMoves := availableMoves(p, g)
	if len(legalMoves) == 0 {
		// Pass turn to opponent
		return minimax(g, depth-1, !isMax, getPlayer(getOpp(p.Piece), g))
	}
	if isMax {
		maxEval := -1000
		for move := range legalMoves {
			gCopy := g
			flipStatic(move, p, gCopy)
			eval := minimax(gCopy, depth-1, false, getPlayer(getOpp(p.Piece), gCopy))
			if eval > maxEval {
				maxEval = eval
			}
		}
		return maxEval
	} else {
		minEval := 1000
		for move := range legalMoves {
			gCopy := g
			flipStatic(move, p, gCopy)
			eval := minimax(gCopy, depth-1, true, getPlayer(getOpp(p.Piece), gCopy))
			if eval < minEval {
				minEval = eval
			}
		}
		return minEval
	}
}

func MaxMove(g *Game, p Player, depth int) int {
	bestMove := -1
	bestValue := -10000
	legalMoves := availableMoves(p, *g)

	for move := range legalMoves {
		gCopy := *g
		flipStatic(move, p, gCopy)
		moveValue := minimax(gCopy, depth, false, getPlayer(getOpp(p.Piece), gCopy))
		if moveValue > bestValue {
			bestValue = moveValue
			bestMove = move
		}
	}
	return bestMove
}

// Get a pseudorandom move for the Randy bot
func RandyMove(moves map[int]bool) int {
	for k := range moves {
		return k
	}
	return 0
}

/*
Represents heuristic values for certain squares on a Board

From: https://courses.cs.washington.edu/courses/cse573/04au/Project/mini1/RUSSIA/Final_Paper.pdf
*/
var StaticWeights = [...]int{
	4, -3, 2, 2, 2, 2, -3, 4, -3, -4, -1, -1, -1, -1, -4, -3, 2, -1, 1, 0, 0,
	1, -1, 2, 2, -1, 0, 1, 1, 0, -1, 2, 2, -1, 0, 1, 1, 0, -1, 2, 2, -1, 1, 0,
	0, 1, -1, 2, -3, -4, -1, -1, -1, -1, -4, -3, 4, -3, 2, 2, 2, 2, -3, 4,
}

func mapToStaticWeights(i int) int {
	row := (i / 10) - 1
	col := (i % 10) - 1
	// fmt.Printf("SW mapped to: %d\n", (row*8 + col))
	return row*8 + col
}
