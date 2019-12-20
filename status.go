package main

import (
	"encoding/json"
	"log"

	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/net"
	packet "github.com/Tnze/go-mc/net/packet"
)

const (
	statusRequest byte = 0x00
	statusPing    byte = 0x01
)

type statusPingPacket struct {
	// May be any number. Notchian clients use a
	// system-dependent time value which is counted
	// in milliseconds.
	Payload packet.Long
}

// statusRequestPacket is used as a server response
type statusRequestPacket struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description chat.Message `json:"description"`
	Favicon     string       `json:"favicon,omitempty"`
}

func handleStatus(conn net.Conn) error {
	for {
		received, err := conn.ReadPacket()
		if err != nil {
			return errConnClosed
		}

		switch received.ID {
		case statusRequest:
			err = handleStatusRequest(conn)
		case statusPing:
			err = handleStatusPing(conn, received)
		}

		if err != nil {
			log.Print("debug: couldn't process status packet (", err, ")")
		}
	}
}

func handleStatusRequest(conn net.Conn) error {
	log.Print("debug: status request packet received")

	status := statusRequestPacket{}
	status.Version.Name = "1.15.1"
	status.Version.Protocol = ProtocolVersion

	status.Players.Max = 100
	status.Players.Online = 0

	status.Description.Text = "pitos de leche"

	jsonResponse, err := json.Marshal(status)
	if err != nil {
		log.Print("error: couldn't marshal server status response")
		return err
	}

	conn.WritePacket(packet.Marshal(
		statusRequest,
		packet.String(string(jsonResponse)),
	))

	return nil
}

func handleStatusPing(conn net.Conn, pingPacket packet.Packet) error {
	log.Print("debug: ping packet received")

	// Parse ping fields.
	ping := statusPingPacket{}
	err := pingPacket.Scan(&ping.Payload)
	if err != nil {
		return errScan
	}

	// Send pong packet. Payload should be the same as sent by the client.
	conn.WritePacket(packet.Marshal(
		statusPing,
		ping.Payload,
	))
	return nil
}
