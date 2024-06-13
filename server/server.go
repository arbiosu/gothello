package server

import (
	"log"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "6578"
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
		go processClient(conn)
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
