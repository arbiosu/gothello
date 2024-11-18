package logic

import (
	"log"
	"strconv"
)

// Used for online play
func (g *Game) GameStatus() (int, int, bool, string, map[int]bool, [100]string) {
	gameOver := g.gameOver()
	legal := g.availableMoves()
	return g.state.score.x, g.state.score.o, gameOver, g.state.currentTurn.piece, legal, g.state.board
}

func (g *Game) MakeOnlineMove(move float64) {
	conv := int(move)
	g.flip(conv)
}

func (g *Game) MakeOnlineBotMove(move int) {
	g.flip(move)
}

// Convert move string to an int
func convertMove(move string) int {
	i, err := strconv.Atoi(move)
	if err != nil {
		log.Println("Error converting move to integer: ", err)
	}
	return i
}
