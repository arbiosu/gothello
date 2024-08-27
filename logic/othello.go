package logic

import (
	"fmt"
)

/* Represents a Player in a Game of Othello */
type Player struct {
	Name  string
	Piece string
}

type PlayerList struct {
	X *Player
	O *Player
}

func NewPlayer(name, piece string) *Player {
	return &Player{
		Name:  name,
		Piece: piece,
	}
}

func NewPlayerList(x, o *Player) *PlayerList {
	return &PlayerList{
		X: x,
		O: o,
	}
}

// Represents a Board in a Game of Othello
type Board struct {
	Board      [100]string
	Turn       string
	Score      *Score
	Players    *PlayerList
	directions [8]int
}

// Print the Board to the console
func (b Board) printBoard() {
	for i := 0; i < len(b.Board); i++ {
		if i%10 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("%v	", b.Board[i])
	}
	fmt.Printf("\n")
}

func newBoard(players *PlayerList) *Board {
	b := [100]string{}
	for i := 0; i < len(b); i++ {
		b[i] = "*"
	}
	for i := 11; i < 89; i++ {
		if validSquare(i) {
			b[i] = "."
		}
	}

	b[44] = "O"
	b[45] = "X"
	b[54] = "X"
	b[55] = "O"

	return &Board{
		Board:      b,
		Turn:       "X",
		Score:      newScore(),
		Players:    players,
		directions: [8]int{10, -10, 1, -1, 11, -11, 9, -9},
	}
}

type Score struct {
	X int
	O int
}

func newScore() *Score {
	return &Score{
		X: 2,
		O: 2,
	}
}

// Captures Board state
type History struct {
	moves []Board
}

func newHistory() *History {
	return &History{}
}

// Represents a Game of Othello. O (white) and X (black) hold their respective
// scores
type Game struct {
	State   *Board
	History *History
}

// Initialize a new game given a PlayerList
func NewGame(players *PlayerList) *Game {
	b := newBoard(players)
	h := newHistory()
	return &Game{
		State:   b,
		History: h,
	}
}

// Returns the current score of a Game. If we are just checking the score,
// does not update the score of the game. TODO: rethink this
func (g *Game) UpdateScore(check bool) (int, int) {
	o, x := 0, 0
	for i := 0; i < len(g.State.Board); i++ {
		switch g.State.Board[i] {
		case "O":
			o += 1
		case "X":
			x += 1
		}
	}
	if !check {
		g.State.Score.O = o
		g.State.Score.X = x
	}
	return o, x
}

// Determines if a square is valid
func validSquare(i int) bool {
	return (i%10 >= 1 && i%10 <= 8)
}

// Returns a string of the opponent's piece
func (g *Game) getOpp() string {
	var opp string
	if g.State.Turn == "X" {
		opp = "O"
	} else {
		opp = "X"
	}
	return opp
}

/* Checks if a given move is valid in a direction */
func (g *Game) validMove(square, direction int) bool {
	// If the square is not an empty square
	if g.State.Board[square] != "." {
		return false
	}

	var (
		opp       = g.getOpp()
		newSquare = square + direction
	)
	// If the next square in that direction is not the opponent's
	if g.State.Board[newSquare] != opp {
		return false
	}
	// Check the rest of the row in that direction
	for validSquare(newSquare) {
		if g.State.Board[newSquare] == "." {
			return false
		}
		// If pieces are sandwiched, return true
		if g.State.Board[newSquare] == g.State.Turn {
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
func (g *Game) validDirections(square int) []int {
	valid := make([]int, 0, 8)
	for i := 0; i < len(g.State.directions); i++ {
		if g.validMove(square, g.State.directions[i]) {
			valid = append(valid, g.State.directions[i])
		}
	}
	return valid
}

func (g *Game) availableMoves() map[int]bool {
	moves := make(map[int]bool)
	for i := 0; i < len(g.State.Board); i++ {
		if validSquare(i) {
			check := g.validDirections(i)
			if len(check) > 0 {
				moves[i] = true
			}
		}
	}
	return moves
}

// Flips the pieces in all valid directions. Records state of the Board.
// Updates Player Turn.
func (g *Game) Flip(square int) {
	opp := g.getOpp()
	if !validSquare(square) {
		if square == -1 {
			g.State.Turn = opp
		}
		return
	}
	dirs := g.validDirections(square)
	g.State.Board[square] = p.Piece

	for i := 0; i < len(dirs); i++ {
		newSquare := square + dirs[i]
		for g.State.Board[newSquare] == opp {
			g.State.Board[newSquare] = p.Piece
			newSquare += dirs[i]
		}
	}
	g.History.moves = append(g.History.moves, *g.State)
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

// Checks if a Game is over
func (g Game) GameOver() bool {
	// Change turn to O
	g.State.Turn = "O"
	o := g.availableMoves()
	// Change turn to X
	g.State.Turn = "X"
	x := g.availableMoves()

	if len(o) < 1 && len(x) < 1 {
		return true
	}

	return false
}

func TestGameOver() {
	p1, p2 := NewPlayer("1", "X"), NewPlayer("2", "O")
	p := NewPlayerList(p1, p2)
	g := NewGame(p)
	test := g.GameOver()
	fmt.Printf("test: %s, turn: %s\n", test, g.State.Turn)
	g.State.Turn = "O"
	test = g.GameOver()
	fmt.Printf("test: %s, turn: %s\n", test, g.State.Turn)
}

func setGameType() int {
	var input int
	fmt.Printf("Enter 1 for to play against a human\n")
	fmt.Printf("Enter 2 to play against Randy, a bot that selects moves at random\n")
	fmt.Printf("Enter 3 to play against Max, a minimax bot\n")
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
