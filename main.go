package main

import (
	"log"

	"github.com/Tnze/go-mc/net"
	"github.com/comail/colog"
)

const (
	// GameVersion is the Minecraft version the server
	// is running for
	GameVersion = "1.15.1"
	// ProtocolVersion is the Minecraft protocol being used
	ProtocolVersion = 575
)

var (
	serverAddress = ""
	serverPort    = "25565"
)

func main() {
	// init logger
	colog.Register()
	colog.SetMinLevel(colog.LDebug)
	colog.SetDefaultLevel(colog.LInfo)

	listener, err := net.ListenMC(serverAddress + ":" + serverPort)
	if err != nil {
		log.Fatal("error: couln't create listener")
		return
	}

	log.Print("info: server running on ", serverAddress, ":", serverPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error: couln't accept client")
			return
		}

		go handleConn(conn)
	}
}
