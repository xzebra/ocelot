package global

import (
	"crypto/rand"
)

var (
	ServerID string
)

func GenerateServerID() {

	rand.Read()
}
