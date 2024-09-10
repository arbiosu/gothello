package server

import (
	"log"
	"net/http"
)

/*
func gameLoop(c *websocket.Conn, data *Data) {
	user := logic.InitializePlayer(data.Name, "X")
	bot := logic.InitializePlayer("Max", "O")
	g, _ := logic.InitializeGame(*user, *bot)
	// Send the board and the legal moves
	_, l, _, _ := logic.OnlineGameStatus(g)
	firstMoves := legalMoves(l)
	data.Legal = firstMoves
	data.Board = g.State.Board
	data.O = 0
	data.X = 0
	d := EncodeData(data)
	c.WriteMessage(websocket.TextMessage, d)
	for !logic.GameOver(*g) {
		curr, _, _, _ := logic.OnlineGameStatus(g)
		if g.State.Turn == "O" {
			time.Sleep(1000 * time.Millisecond)
			move := logic.MaxMove(g, curr, 3)
			logic.Flip(move, curr, g)
			_, userMoves, o, x := logic.OnlineGameStatus(g)
			data.Legal = legalMoves(userMoves)
			data.Board = g.State.Board
			data.O = o
			data.X = x
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
			move := convertMove(game.Move)
			logic.Flip(move, curr, g)
			_, _, o, x := logic.OnlineGameStatus(g)
			game.Board = g.State.Board
			game.O = o
			game.X = x
			encoded := EncodeData(&game)
			c.WriteMessage(msgType, encoded)
		}
	}
	log.Println("GAME OVER")
	data.GameOver = 1
	encoded := EncodeData(data)
	c.WriteMessage(websocket.TextMessage, encoded)
}
*/

func Server() {
	hub := NewHub()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.handleWs)
	s := &http.Server{
		Addr:    ":7342",
		Handler: mux,
	}

	log.Println("Initializing server...Go!")
	log.Fatal(s.ListenAndServe())
}
