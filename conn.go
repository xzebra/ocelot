package main

import (
	"log"

	"github.com/Tnze/go-mc/net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	info, err := handleHandshake(conn)
	if err != nil {
		log.Print("debug: couldn't get handshake (", err, ")")
		return
	}

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
