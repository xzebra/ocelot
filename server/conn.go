package server

import (
	"errors"
	"log"

	"github.com/Tnze/go-mc/net"
)

var (
	errConnClosed = errors.New("connection closed")
	errScan       = errors.New("couldn't parse packet fields")
)

// HandleConn handles a client connection
func HandleConn(conn net.Conn) {
	defer conn.Close()

	log.Print("info: client stablished connection")

	info, err := handleHandshake(conn)
	if err != nil {
		log.Print("debug: couldn't get handshake (", err, ")")
		return
	}

	log.Print("debug: handshake successful")

	switch info.NextState {
	case nextStateStatus:
		handleStatus(conn)
		// once status has been handled, close connection
		return
	case nextStateLogin:
		err = handleLogin(conn)
		if err != nil {
			log.Print("debug: player couldn't login (", err, ")")
			return
		}
		err = handlePlay(conn)
		if err != nil {
			log.Print("debug: error while playing (", err, ")")
			return
		}
	}
}
