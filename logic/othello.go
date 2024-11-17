package logic

func NewPlayer(n, p string) *Player {
	return &Player{
		name:  n,
		piece: p,
	}
}

// Always put the X player first TODO: rethink
func NewPlayerList(x, o *Player) *PlayerList {
	return &PlayerList{
		x: x,
		o: o,
	}
}

// Initialize a new board
// TODO: one loop
func newBoard(players *PlayerList) *GameData {
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

	return &GameData{
		board:       b,
		currentTurn: players.x,
		score:       newScore(),
		players:     players,
		directions:  [8]int{10, -10, 1, -1, 11, -11, 9, -9},
	}
}

func newScore() *score {
	return &score{
		x: 2,
		o: 2,
	}
}

func newHistory() *history {
	return &history{}
}

// Initialize a new game given a PlayerList
func NewGame(players *PlayerList) *Game {
	b := newBoard(players)
	h := newHistory()
	return &Game{
		state:   b,
		history: h,
	}
}

// Updates the current score of a Game.
func (g *Game) UpdateScore() {
	o, x := 0, 0
	for i := 0; i < len(g.state.board); i++ {
		switch g.state.board[i] {
		case "O":
			o += 1
		case "X":
			x += 1
		}
	}
	g.state.score.o = o
	g.state.score.x = x
}

// Get the opposing player for the current player
func (g *Game) getOpponent() *Player {
	var opp *Player
	switch g.state.currentTurn.piece {
	case "O":
		opp = g.state.players.x
	case "X":
		opp = g.state.players.o
	}
	return opp
}

// Checks if a move is valid for the current player's turn
func (g *Game) validMove(square, direction int) bool {
	// If the initial square is not an empty square, its not a valid move
	if g.state.board[square] != "." {
		return false
	}

	newSquare := square + direction
	opp := g.getOpponent()

	if g.state.board[newSquare] != opp.piece {
		return false
	}

	for validSquare(newSquare) {
		if g.state.board[newSquare] == "." || g.state.board[newSquare] == "*" {
			return false
		}
		// If pieces are sandwiched, its a valid move!
		if g.state.board[newSquare] == g.state.currentTurn.piece {
			return true
		}
		newSquare += direction
	}
	return false
}

// Returns a slice of all valid directions for a given move
func (g *Game) validDirections(square int) []int {
	valid := make([]int, 0, 8)
	for i := 0; i < len(g.state.directions); i++ {
		if g.validMove(square, g.state.directions[i]) {
			valid = append(valid, g.state.directions[i])
		}
	}
	return valid
}

// Returns a map of the current player's available moves
func (g *Game) availableMoves() map[int]bool {
	moves := make(map[int]bool)
	for i := 0; i < len(g.state.board); i++ {
		if validSquare(i) {
			check := g.validDirections(i)
			if len(check) > 0 {
				moves[i] = true
			}
		}
	}
	return moves
}

// Flips the pieces in all valid directions for the given square and player
// Records the current state into a Game's history
// Updates the score
func (g *Game) flip(square int) {
	opp := g.getOpponent()
	// if -1 is given, pass the turn to the other player
	// symbolizes when the current player has no available moves
	if !validSquare(square) {
		if square == -1 {
			g.state.currentTurn = opp
		}
		return
	}
	g.history.moves = append(g.history.moves, *g.state)
	directions := g.validDirections(square)
	g.state.board[square] = g.state.currentTurn.piece
	// Flip the pieces in every valid direction until we reach the sandwiched
	// piece
	for i := 0; i < len(directions); i++ {
		newSquare := square + directions[i]
		for g.state.board[newSquare] == opp.piece {
			g.state.board[newSquare] = g.state.currentTurn.piece
			newSquare += directions[i]
		}
	}
	g.UpdateScore()
	g.state.currentTurn = opp
}

// Checks if a game is over
func (g *Game) gameOver() bool {
	// Save the current player's turn
	current := g.state.currentTurn
	opp := g.getOpponent()

	// get available moves for current player
	currentMoves := g.availableMoves()
	// change turn to opponent and get opp moves
	g.state.currentTurn = opp
	oppMoves := g.availableMoves()
	// revert changes
	g.state.currentTurn = current
	// The game is over when both players have no available moves
	if len(currentMoves) < 1 && len(oppMoves) < 1 {
		return true
	}
	return false
}
