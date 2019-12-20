package main

import (
	"log"

	"github.com/Tnze/go-mc/net"
)

func HandleConn(conn net.Conn) {
	for {
		packet, err := conn.ReadPacket()
		if err != nil {
			log.Print("info: could't read from player")
			break
		}

		log.Print("debug: received packet with id = ", packet.ID)
	}

	conn.Close()
}

func HandleLogin() error {
	// Login process (C = client, S = server)
	// C→S: Handshake with Next State set to 2 (login)

	// C→S: Login Start
	// For unauthenticated and localhost connections there is no encryption.
	// S→C: Encryption Request
	// Client auth
	// C→S: Encryption Response
	// Server auth, both enable encryption
	// S→C: Set Compression (optional)
	// S→C: Login Success

	return nil
}
