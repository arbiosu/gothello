package logic

import (
	"fmt"
	"os"
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
// TODO: make it return a Player
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

// Returns a slice of all valid directions for a given move
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
	g.State.Board[square] = g.State.Turn

	for i := 0; i < len(dirs); i++ {
		newSquare := square + dirs[i]
		for g.State.Board[newSquare] == opp {
			g.State.Board[newSquare] = g.State.Turn
			newSquare += dirs[i]
		}
	}
	g.History.moves = append(g.History.moves, *g.State)
	g.State.Turn = opp
}

// Used for bot play and calculation - does not affect actual state
func flipStatic(square int, g Game) {
	opp := g.getOpp()
	dirs := g.validDirections(square)
	g.State.Board[square] = g.State.Turn

	for i := 0; i < len(dirs); i++ {
		newSquare := square + dirs[i]
		for g.State.Board[newSquare] == opp {
			g.State.Board[newSquare] = g.State.Turn
			newSquare += dirs[i]
		}
	}
	g.State.Turn = opp
}

// Checks if a Game is over
// TODO: rethink this and availableMoves func because i dont like this
func (g *Game) GameOver() bool {
	curr := g.getCurrentPlayer().Piece
	// Change turn to O
	g.State.Turn = "O"
	o := g.availableMoves()
	// Change turn to X
	g.State.Turn = "X"
	x := g.availableMoves()

	g.State.Turn = curr

	if len(o) < 1 && len(x) < 1 {
		return true
	}

	return false
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

func (g *Game) getCurrentPlayer() *Player {
	var curr *Player
	if g.State.Turn == "O" {
		curr = g.State.Players.O
	} else {
		curr = g.State.Players.X
	}
	return curr
}

// Used for minimax evaluation as Turn is not tracked. TODO: rethink
func (g *Game) getOppPlayer(curr Player) *Player {
	var opp *Player
	if curr.Piece == "O" {
		opp = g.State.Players.X
	} else {
		opp = g.State.Players.O
	}
	return opp
}

// TODO: rethink
func result(g *Game) {
	o, x := g.UpdateScore(false)
	fmt.Println("------GAME OVER!------")
	g.State.printBoard()
	fmt.Printf("FINAL SCORE:	O: %d	X: %d\n", o, x)
	switch {
	case o > x:
		winner := g.State.Players.O
		fmt.Printf("WINNER: %v. CONGRATULATIONS!\n", winner.Name)
	case o < x:
		winner := g.State.Players.X
		fmt.Printf("WINNER: %v. CONGRATULATIONS!\n", winner.Name)
	default:
		fmt.Printf("RESULT: TIE GAME!\n")
	}
	// fmt.Println("Thank you for playing Gothello!")
}

// Updates the current status of a Game. Returns a map of the available moves for
// the Player whose turn it is.
func (g *Game) GameStatus() map[int]bool {
	var (
		o, x  = g.UpdateScore(true)
		curr  = g.getCurrentPlayer()
		legal = g.availableMoves()
	)

	fmt.Printf("CURRENT TURN: %v, %s\n", curr.Name, g.State.Turn)
	fmt.Printf("SCORE:	O: %d X: %d\n", o, x)
	fmt.Printf("LEGAL MOVES FOR %s: %v\n", g.State.Turn, legal)
	return legal
}

// Process a Human Move in the terminal
func (g *Game) humanMove(legal map[int]bool) {
	var move int
	curr := g.getCurrentPlayer()
	fmt.Printf("%v, enter your move. Choose any of the above legal moves: ", curr.Name)
	fmt.Scanln(&move)
	elem, ok := legal[move]
	if !ok {
		fmt.Printf("%v INVALID MOVE!\n", elem)
	} else {
		g.Flip(move)
	}
}

// Set up the players for a game against two Humans in the terminal.
func setPlayersHuman() *PlayerList {
	var name, name2 string
	fmt.Printf("Enter Player 1's name. They will be the \"O\" pieces: ")
	fmt.Scanln(&name)
	fmt.Printf("Enter Player 2's name. They will be the \"X\" pieces: ")
	fmt.Scanln(&name2)
	p1, p2 := NewPlayer(name, "O"), NewPlayer(name2, "X")
	players := NewPlayerList(p2, p1)
	return players
}

// Set up the players for a game against a Human and a Bot in the terminal.
func setPlayersBot(bot string) *PlayerList {
	var name string
	fmt.Printf("Initializing %s...complete! %s will play as `O`\n", bot, bot)
	fmt.Printf("Enter your name: ")
	fmt.Scanln(&name)
	fmt.Printf("\nWelcome, %s. Good luck against %s!\n", name, bot)
	p1, p2 := NewPlayer(bot, "O"), NewPlayer(name, "X")
	players := NewPlayerList(p2, p1)
	return players
}

func HumanGame() {
	p := setPlayersHuman()
	g := NewGame(p)
	fmt.Printf("PLAYERS: (%s: %s), (%s: %s)\n", p.O.Name, p.O.Piece, p.X.Name, p.X.Piece)
	for !g.GameOver() {
		g.State.printBoard()
		legal := g.GameStatus()
		g.humanMove(legal)
	}
	result(g)
}

func BotGame(bot string) {
	p := setPlayersBot(bot)
	g := NewGame(p)
	fmt.Printf("PLAYERS: (%s: %s), (%s: %s)\n", p.O.Name, p.O.Piece, p.X.Name, p.X.Piece)
	for !g.GameOver() {
		g.State.printBoard()
		legal := g.GameStatus()
		var move int
		if g.State.Turn == "O" {
			fmt.Printf("%s is thinking on a move...\n", bot)
			switch bot {
			case "Randy":
				move = RandyMove(legal)
			case "Max":
				move = MaxMove(g, 3)
			}
			fmt.Printf("%s chose: %d\n", move)
			g.Flip(move)
		} else {
			g.humanMove(legal)
		}
	}
	result(g)
}

func displayHelpMsg() {
	fmt.Println("TODO: PRINT OTHELLO RULES")
}

// Play Othello in the terminal!
func PlayGame() {
	gameType := setGameType()
	switch gameType {
	case 1:
		HumanGame()
	case 2:
		BotGame("Randy")
	case 3:
		BotGame("Max")
	case 4:
		displayHelpMsg()
	default:
		fmt.Println("EXITING....")
	}
	os.Exit(0)
}

// TODO: Test and refactor minimax portion
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

func (g *Game) minimax(depth int, isMax bool, p Player) int {
	if g.GameOver() || depth == 0 {
		return heuristic(g, p)
	}
	legal := g.availableMoves()
	if len(legal) == 0 {
		// Pass minimax evaluation turn to opponent
		return g.minimax(depth-1, !isMax, *g.getOppPlayer(p))
	}
	if isMax {
		maxEval := -1000
		for move := range legal {
			gCopy := *g
			flipStatic(move, gCopy)
			eval := gCopy.minimax(depth-1, false, *g.getOppPlayer(p))
			if eval > maxEval {
				maxEval = eval
			}
		}
		return maxEval
	} else {
		minEval := 1000
		for move := range legal {
			gCopy := *g
			flipStatic(move, gCopy)
			eval := g.minimax(depth-1, true, *g.getOppPlayer(p))
			if eval < minEval {
				minEval = eval
			}
		}
		return minEval
	}
}

func MaxMove(g *Game, depth int) int {
	bestMove := -1
	bestValue := -10000
	legal := g.availableMoves()

	for move := range legal {
		gCopy := *g
		flipStatic(move, gCopy)
		// Bots are always "O"
		moveValue := gCopy.minimax(depth, false, *g.State.Players.X)
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
