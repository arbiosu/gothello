package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type OnlinePlayer struct {
	Name string `json:"name"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Might need to change
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}

func socketReader(conn *websocket.Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("socketReader: Error reading message from client! (%v)", err)
			break
		}
		log.Printf("Received: %s", msg)
		var user OnlinePlayer
		err = json.Unmarshal(msg, &user)
		if err != nil {
			log.Printf("socketReader: Error Unmarshaling JSON! (%v)", err)
		}
		err = conn.WriteMessage(msgType, msg)
		if err != nil {
			log.Printf("socketReader: Error writing message! (%v)", err)
		}
	}
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected")
	socketReader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWs)
}

func Server() {
	setupRoutes()
	log.Println("Initializing server...Go!")
	log.Fatal(http.ListenAndServe(":7333", nil))
}
