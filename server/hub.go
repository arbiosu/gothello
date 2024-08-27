package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Hub struct {
	clients ClientList
	sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(ClientList),
	}
}

func (h *Hub) handleWs(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
	w.Header().Set("Content-Type", "application/json")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, h)

	h.addClient(client)

	go client.ReadMessages()

}

func (h *Hub) addClient(c *Client) {
	h.Lock()
	defer h.Unlock()

	h.clients[c] = true
}

func (h *Hub) removeClient(c *Client) {
	h.Lock()
	defer h.Unlock()

	// Check if client is in the Hub and delete it.
	if _, ok := h.clients[c]; ok {
		c.conn.Close()
		delete(h.clients, c)
	}
}
