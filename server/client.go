package server

import (
	"log"

	"github.com/gorilla/websocket"
)

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
