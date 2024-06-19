package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/gothello/bots"
	"github.com/gothello/logic"
)

// Represents the data client and server will send back and forth

type Data struct {
	Name  string      `json:"name"`
	Board [100]string `json:"board"`
	Move  string      `json:"move"`
	Legal []int       `json:"legal"`
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
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("socketReader: Error reading message from client! (%v)", err)
			break
		}
		log.Printf("sr: Received: %s", msg)
		var game Data
		err = json.Unmarshal(msg, &game)
		if err != nil {
			log.Printf("socketReader: Error Unmarshaling JSON! (%v)", err)
		}
		gameLoop(conn, &game)
	}
	closing := "Thanks for playing!"
	msg, err := json.Marshal(closing)
	if err != nil {
		log.Println("Error encoding closing msg: ", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("Error writing closing msg: ", err)
	}
}

func EncodeData(d *Data) []byte {
	res, err := json.Marshal(d)
	if err != nil {
		log.Println("ERROR encoding state: ", err)
	}
	return res
}

func legalMoves(m map[int]bool) []int {
	var legal []int
	for k, _ := range m {
		legal = append(legal, k)
	}
	return legal
}

func gameLoop(c *websocket.Conn, data *Data) {
	user := logic.InitializePlayer(data.Name, "X")
	bot := logic.InitializePlayer("Randy", "O")
	g, _ := logic.InitializeGame(*user, *bot)
	// Send the board and the legal moves
	_, l := logic.OnlineGameStatus(g)
	firstMoves := legalMoves(l)
	data.Legal = firstMoves
	data.Board = g.State.Board
	d := EncodeData(data)
	c.WriteMessage(websocket.TextMessage, d)
	for !logic.GameOver(*g) {
		curr, legal := logic.OnlineGameStatus(g)
		if g.State.Turn == "O" {
			move := bots.RandyMove(legal)
			logic.Flip(move, curr, g)
			_, userMoves := logic.OnlineGameStatus(g)
			data.Legal = legalMoves(userMoves)
			data.Board = g.State.Board
			encoded := EncodeData(data)
			c.WriteMessage(websocket.TextMessage, encoded)
		} else {
			// Read message
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("gameLoop: Error reading message from client! (%v)", err)
				break
			}
			// Decode
			var game Data
			err = json.Unmarshal(msg, &game)
			if err != nil {
				log.Printf("gameLoop: Error Unmarshaling JSON! (%v)", err)
			}
			// Perform backend logic: make move, send board back
			move := convertMove(data.Move)
			logic.Flip(move, curr, g)
			encoded := EncodeData(&game)
			c.WriteMessage(msgType, encoded)
		}
	}
}

func convertMove(move string) int {
	i, err := strconv.Atoi(move)
	if err != nil {
		log.Println("Error converting move to integer: ", err)
	}
	return i
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
	log.Fatal(http.ListenAndServe(":7340", nil))
}
