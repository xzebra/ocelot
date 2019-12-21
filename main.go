package main

import (
	"log"

	"ocelot/global"
	"ocelot/server"

	"github.com/Tnze/go-mc/net"
	"github.com/comail/colog"
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

	global.GenerateKeyPair()

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

		go server.HandleConn(conn)
	}
}
