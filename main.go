package main

import (
	"log"
	"net/http"

	"github.com/gothello/api"
	"github.com/gothello/logic"
)

func main() {
	// logic.PlayGame()
	http.HandleFunc("/test", test)
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
	res := api.EncodeState(bptr)
	w.Write(res)
}
