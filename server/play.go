package server

import (
	"log"

	"github.com/Tnze/go-mc/net"
)

func handlePlay(conn net.Conn) error {
	for {
		packet, err := conn.ReadPacket()
		if err != nil {
			log.Print("debug: could't read from player")
			return errConnClosed
		}

		log.Print("debug: received packet with id = ", packet.ID)
	}
}
