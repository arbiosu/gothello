package server

import "github.com/gorilla/websocket"

type ClientList map[*Client]bool

type Client struct {
	conn *websocket.Conn

	hub *Hub
}

func NewClient(c *websocket.Conn, h *Hub) *Client {
	return &Client{
		conn: c,
		hub:  h,
	}
}
