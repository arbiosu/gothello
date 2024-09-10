package server

import (
	"log"

	"github.com/gorilla/websocket"
)

// Represents the data client and server will send back and forth
type GameData struct {
	Name     string      `json:"name"`
	Opp      string      `jaon:"opp"` // bot opponent name
	Board    [100]string `json:"board"`
	Move     string      `json:"move"`
	Legal    []int       `json:"legal"`
	O        int         `json:"o"`
	X        int         `json:"x"`
	GameOver int         `json:"gameOver"`
}

type ClientList map[*Client]bool

type Client struct {
	conn *websocket.Conn
	hub  *Hub

	// ch is used to avoid concurrent writes
	ch chan []byte
}

func NewClient(c *websocket.Conn, h *Hub) *Client {
	return &Client{
		conn: c,
		hub:  h,
		ch:   make(chan []byte),
	}
}

// Read a message from a Client
func (c *Client) readMessages() {
	defer func() {
		// Cleanup connection
		c.hub.removeClient(c)
	}()

	for {
		msgType, p, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: could not read message from client! (%v)\n", err)
				// TODO: handle
			}
			break
		}
		log.Println(msgType)
		log.Println(p)
	}
}

// Write a message to a client
func (c *Client) writeMessages() {
	defer func() {
		c.hub.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.ch:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("Error: connection closed! (%v)\n", err)
				}
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error: failed to send message (%v)\n", err)
			}
			log.Println("Message sent")
		}
	}
}
