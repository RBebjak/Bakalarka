package main

import (
	"goChat/server/clientHandle"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Println("TCP chat server running on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go clientHandle.HandleConnection(conn)
	}
}
