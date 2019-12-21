package global

import (
	"log"

	"crypto/rand"
	"crypto/rsa"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

// GenerateKeyPair generates a 1024-bit RSA keypair
func GenerateKeyPair() {
	log.Print("debug: generating RSA keypair")

	var err error
	PrivateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal("error: couldn't generate encryption keys")
	}

	PublicKey = &PrivateKey.PublicKey
}
