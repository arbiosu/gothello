package logic

import (
	"log"
	"strconv"
)

// Used for online play
func (g *Game) GameStatus() (*score, bool, *Player, map[int]bool, [100]string) {
	gameOver := g.gameOver()
	legal := g.availableMoves()
	return g.state.score, gameOver, g.state.currentTurn, legal, g.state.board
}

func (g *Game) MakeOnlineMove(move string) {
	conversion := convertMove(move)
	g.flip(conversion)
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
