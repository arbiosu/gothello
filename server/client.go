package server

import (
	"log"
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
			move := msg.Content.(float64)
			c.Lock()
			c.game.MakeOnlineMove(move)
			log.Printf("client making move %f", move)
			c.sendGameState()
			c.Unlock()
			// Wait for a bit before sending bot response move
			time.Sleep(1000 * time.Millisecond)
			botMove := logic.MaxMove(c.game, 3)
			c.Lock()
			c.game.MakeOnlineBotMove(botMove)
			c.sendGameState()
			c.Unlock()
		case "newGame":
			// TODO FIX
			c.game = logic.NewGame(logic.NewPlayerList(logic.NewPlayer("Human", "X"), logic.NewPlayer("Max", "O")))
			c.sendGameState()
		}
	}
}

func (c *Client) sendGameState() {
	x, o, gameOver, turn, legal, board := c.game.GameStatus()

	state := Message{
		Type: "gameState",
		Content: map[string]interface{}{
			"board":    board,
			"turn":     turn,
			"legal":    legal,
			"x":        x,
			"o":        o,
			"gameOver": gameOver,
		},
	}
	c.Lock()
	err := c.conn.WriteJSON(state)
	c.Unlock()
	if err != nil {
		log.Printf("Err: %v\n", err)
		// TODO: handle
	}
}
