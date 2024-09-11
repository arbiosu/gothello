package server

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gothello/logic"
)

// Represents the messages client and server will send back and forth
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type ClientList map[*Client]bool

type Client struct {
	conn *websocket.Conn
	hub  *Hub
	game *logic.Game
	sync.RWMutex
}

func NewClient(c *websocket.Conn, h *Hub, pl *logic.PlayerList) *Client {
	return &Client{
		conn: c,
		hub:  h,
		game: logic.NewGame(pl),
	}
}

func (c *Client) playGame() {
	defer func() {
		// Cleanup connection
		c.hub.removeClient(c)
	}()

	// Send initial game state
	c.sendGameState()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Err: %v\n", err)
			break
		}
		switch msg.Type {
		case "move":
			move := msg.Content.(string)
			conv := convertMove(move)
			// TODO: mutex
			c.game.Flip(conv)
			c.sendGameState()
			if c.game.State.Turn == "O" {
				// Wait for a bit before sending next move
				// TODO: mutex
				time.Sleep(1000 * time.Millisecond)
				aiMove := logic.MaxMove(c.game, 3)
				c.game.Flip(aiMove)
				c.sendGameState()
			}
		case "newGame":
			// TODO
			c.game = logic.NewGame(logic.NewPlayerList(logic.NewPlayer("Human", "X"), logic.NewPlayer("Max", "O")))
			c.sendGameState()
		}
	}
}

func (c *Client) sendGameState() {
	o, x := c.game.UpdateScore(false)
	board := c.game.State.Board
	turn := c.game.State.Turn
	gameOver := c.game.GameOver()
	legal := c.game.AvailableMoves()

	state := Message{
		Type: "gameState",
		Content: map[string]interface{}{
			"board":    board,
			"turn":     turn,
			"legal":    legal,
			"o":        o,
			"x":        x,
			"gameOver": gameOver,
		},
	}

	err := c.conn.WriteJSON(state)
	if err != nil {
		log.Printf("Err: %v\n", err)
		// TODO: handle
	}
}

// Convert move string to an int TODO: rethink on js side
func convertMove(move string) int {
	i, err := strconv.Atoi(move)
	if err != nil {
		log.Println("Error converting move to integer: ", err)
	}
	return i
}
