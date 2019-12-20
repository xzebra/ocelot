package server

import (
	"errors"
	"log"

	"github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/net/packet"
)

const (
	// Handshake packet ID
	Handshake byte = 0x00
)

const (
	nextStateError packet.VarInt = iota
	nextStateStatus
	nextStateLogin
)

var (
	errHandshakeFormat = errors.New("wrong handshake packet format")
)

type handshakePacket struct {
	// ProtocolVersion is the version that the client plans
	// on using to connect to the server (which is not
	// important for the ping).
	// If the client is pinging to determine what version to
	// use, by convention -1 should be set.
	ProtocolVersion packet.VarInt

	// Hostname or IP, e.g. localhost or 127.0.0.1, that
	// was used to connect
	ServerAddress packet.String

	ServerPort packet.UnsignedShort

	// NextState should be 1 for status, but could also be 2 for login.
	NextState packet.VarInt
}

// handleHandshake returns the next server state (status
// or login) or an error
func handleHandshake(conn net.Conn) (handshakePacket, error) {
	hp := handshakePacket{}

	for {
		packet, err := conn.ReadPacket()
		if err != nil {
			return hp, errConnClosed
		}

		err = packet.Scan(&hp.ProtocolVersion,
			&hp.ServerAddress, &hp.ServerPort,
			&hp.NextState)
		if err != nil {
			log.Print("debug: unknown handshake format")
		} else {
			break
		}
	}

	return hp, nil
}
