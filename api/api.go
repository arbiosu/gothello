package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gothello/logic"
)

/* POST structs */
type Move struct {
	Square int `json:"move"`
}

type OnlinePlayer struct {
	Name  string `json:"name"`
	Piece string `json:"piece"`
}

func GoServer() {
	http.HandleFunc("/test", test)
	http.HandleFunc("/init", initOnlinePlayer)
	log.Println("Go!")
	http.ListenAndServe(":8080", nil)
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	p1 := logic.InitializePlayer("Arbi", "O")
	p2 := logic.InitializePlayer("Tayler", "X")
	_, bptr := logic.InitializeGame(*p1, *p2)
	res := encodeState(bptr)
	w.Write(res)
}

func initOnlinePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "POST":
		var user OnlinePlayer
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		log.Printf("%s, %s", user.Name, user.Piece)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "ERR: not a POST request")
	}
}

func showState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func encodeState(b *logic.Board) []byte {
	res, err := json.Marshal(b.Board)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return res
}
