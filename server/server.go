package server

import (
	"log"
	"net"
	"os"

	"github.com/gothello/logic"
)

const (
	HOST = "localhost"
	PORT = "8080"
)

func GoServer() {
	log.Println("Establishing tcp server...")
	server, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		log.Printf("FAILED! Error: %v", err.Error())
	}
	defer server.Close()
	log.Printf("Listening on %s:%s", HOST, PORT)
	log.Println("Waiting for client...")
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err.Error())
			os.Exit(1)
		}
		log.Println("Accepted client!")
		go startGame(conn)
	}
}

func processClient(conn net.Conn) {
	buf := make([]byte, 1024)
	msg, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error: %v", err.Error())
	}
	log.Printf("MSG: %s", string(buf[:msg]))
	_, err = conn.Write([]byte("MSG RECEIVED"))
	if err != nil {
		log.Printf("Error writing message: %v", err.Error())
	}
	conn.Close()
}

func startGame(conn net.Conn) {
	intro := []byte("Enter your name: ")
	_, err := conn.Write(intro)
	if err != nil {
		log.Printf("Error writing message: %v", err.Error())
	}
	buf := make([]byte, 1024)
	msg, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error reading message: %v", err.Error())
	}
	p := &logic.Player{Name: string(msg), Piece: "X"}
	randy := &logic.Player{Name: "Randy", Piece: "O"}
	gptr, _ := logic.InitializeGame(*p, *randy)
	encoded := logic.EncodeState(gptr)
	_, err = conn.Write(encoded)
	if err != nil {
		log.Printf("Error writing state: %v", err.Error())
	}
}
