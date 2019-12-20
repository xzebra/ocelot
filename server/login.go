package server

import (
	"github.com/Tnze/go-mc/net"
)

func handleLogin(conn net.Conn) error {
	// Login process (C = client, S = server):

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
