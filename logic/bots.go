package logic

import "math"

// Get a pseudorandom move for the Randy bot
func RandyMove(moves map[int]bool) int {
	for k := range moves {
		return k
	}
	return 0
}

// Get a move chosen by the minimax algorithm
func MaxMove(g *Game, depth int) int {
	bestMove := -1
	bestScore := math.Inf(-1)

	for move := range g.availableMoves() {
		gameCopy := copyGame(g)
		gameCopy.flip(move)

		score := minimax(gameCopy, depth-1, false)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return bestMove
}

// Minimax alogrithm
func minimax(g *Game, depth int, maximizingPlayer bool) float64 {
	if depth == 0 || g.gameOver() {
		return evaluateBoard(g)
	}

	if maximizingPlayer {
		maxEval := math.Inf(-1)
		for move := range g.availableMoves() {
			gameCopy := copyGame(g)
			gameCopy.flip(move)
			eval := minimax(gameCopy, depth-1, false)
			maxEval = math.Max(maxEval, eval)
		}
		return maxEval
	} else {
		minEval := math.Inf(1)
		for move := range g.availableMoves() {
			gameCopy := copyGame(g)
			gameCopy.flip(move)
			eval := minimax(gameCopy, depth-1, true)
			minEval = math.Min(minEval, eval)
		}
		return minEval
	}
}

// Assigns a score to the current board state
// TODO: fix
func evaluateBoard(g *Game) float64 {
	g.UpdateScore()
	return float64(g.state.score.o)
}

// Creates a deep copy of the game state
func copyGame(g *Game) *Game {
	newBoard := *g.state
	newBoard.board = g.state.board
	newGame := &Game{
		state:   &newBoard,
		history: g.history,
	}
	return newGame
}

// Determines if a square is valid
func validSquare(i int) bool {
	return (i%10 >= 1 && i%10 <= 8)
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

/*
func mapToStaticWeights(i int) int {
	row := (i / 10) - 1
	col := (i % 10) - 1
	// fmt.Printf("SW mapped to: %d\n", (row*8 + col))
	return row*8 + col
}
*/
