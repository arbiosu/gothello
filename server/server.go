package server

import (
	"log"
	"net/http"
)

func Server() {
	hub := NewHub()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.handleWs)
	s := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	log.Println("Initializing server...Go!")
	log.Fatal(s.ListenAndServe())
}
