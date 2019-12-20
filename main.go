package main

import (
	"log"

	"github.com/Tnze/go-mc/net"
	"github.com/comail/colog"
)

func main() {
	// init logger
	colog.Register()
	colog.SetMinLevel(colog.LDebug)
	colog.SetDefaultLevel(colog.LInfo)

	listener, err := net.ListenMC(":25565")
	if err != nil {
		log.Fatal("error: couln't create listener")
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error: couln't accept client")
			return
		}

		go HandleConn(conn)
	}
}
