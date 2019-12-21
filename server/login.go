package server

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"

	"ocelot/global"

	"github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/net/CFB8"
	"github.com/Tnze/go-mc/net/packet"
)

const (
	verifyTokenLength = 4
)

const (
	// Client packets
	LoginStart          byte = 0x00
	EncryptionResponse  byte = 0x01
	LoginPluginResponse byte = 0x02

	// Server packets
	Disconnect         byte = 0x00
	EncryptionRequest  byte = 0x01
	LoginSuccess       byte = 0x02
	SetCompression     byte = 0x03
	LoginPluginRequest byte = 0x04
)

var (
	errNoLoginStart         = errors.New("no login start packet")
	errNoEncryptionResponse = errors.New("no encryption response packet")
	errPublicKeyMarshal     = errors.New("couldn't marshal public key in ASN.1 format")
	errWrongSecretLength    = errors.New("secret length doesn't match received length")
	errWrongTokenLength     = errors.New("verify token length doesn't match received length")
	errVerifyToken          = errors.New("verify token doesn't match")
)

type encryptionRequestPacket struct {
	// ServerID appears to be empty
	ServerID packet.String
	// PublicKeyLength is the length of PublicKey
	PublicKeyLength packet.VarInt
	// PublicKey
	PublicKey packet.ByteArray
	// VerifyTokenLength Length of VerifyToken.
	// Always 4 for Notchian servers.
	VerifyTokenLength packet.VarInt
	// VerifyToken is a sequence of random bytes
	// generated by the server.
	VerifyToken packet.ByteArray
}

func (p *encryptionRequestPacket) Marshal() (pk packet.Packet) {
	return packet.Marshal(
		EncryptionRequest,
		p.ServerID,
		p.PublicKey,
		p.VerifyToken,
	)
}

type encryptionResponsePacket struct {
	// SecretLength stores the length of shared secret
	SecretLength packet.VarInt
	// Secret is the shared secret between client and server
	Secret packet.ByteArray
	// VerifyTokenLength Length of VerifyToken.
	// Should be the same as sent by server.
	VerifyTokenLength packet.VarInt
	// VerifyToken is a sequence of random bytes
	// generated by the server, but this one is
	// encrypted with server's public key.
	VerifyToken packet.ByteArray
}

func (p *encryptionResponsePacket) Decode(r packet.DecodeReader) error {
	if err := p.Secret.Decode(r); err != nil {
		return errScan
	}
	p.SecretLength = packet.VarInt(len(p.Secret))

	if err := p.VerifyToken.Decode(r); err != nil {
		return errScan
	}
	p.VerifyTokenLength = packet.VarInt(len(p.VerifyToken))

	return nil
}

type loginSuccessPacket struct {
	// UUID unlike in other packets, this field
	// contains the UUID as a string with hyphens.
	UUID     packet.String
	Username packet.String
}

func (p *loginSuccessPacket) Marshal() (pk packet.Packet) {
	return packet.Marshal(
		LoginSuccess,
		p.UUID,
		p.Username,
	)
}

func handleLogin(conn net.Conn) error {
	// Login process (C = client, S = server):

	// C→S: Login Start
	username, err := handleLoginStart(conn)
	if err != nil {
		return err
	}
	// TODO: For unauthenticated and localhost connections there is no encryption.

	// S→C: Encryption Request
	verifyToken, err := handleEncryptionRequest(conn)
	if err != nil {
		return err
	}

	secret, err := handleEncryptionResponse(conn, verifyToken)

	// Server auth
	uuid, err := handleServerAuth()

	// Both enable AES/CFB8 encryption using the shared secret as both IV and key.
	block, err := aes.NewCipher(secret)
	if err != nil {
		return err
	}
	conn.SetCipher(
		CFB8.NewCFB8Encrypt(block, secret),
		CFB8.NewCFB8Decrypt(block, secret),
	)

	// TODO: S→C: Set Compression (optional)

	// S→C: Login Success
	return handleLoginSuccess(conn, username, uuid)
}

func handleLoginStart(conn net.Conn) (packet.String, error) {
	var username packet.String

	received, err := conn.ReadPacket()
	if err != nil {
		return username, errConnClosed
	}
	// Check if it really is a LoginStart packet
	if received.ID != LoginStart {
		return username, errNoLoginStart
	}

	// Scan username parameter
	err = received.Scan(&username)
	if err != nil {
		return username, errScan
	}

	// TODO: check username format...

	return username, nil
}

func handleEncryptionRequest(conn net.Conn) ([]byte, error) {
	// Generate random verify token
	verifyToken := make([]byte, verifyTokenLength)
	rand.Read(verifyToken)

	// We have to send the public key in PKIX, ASN.1 DER format.
	asn1PublicKey, err := x509.MarshalPKIXPublicKey(global.PublicKey)
	if err != nil {
		return verifyToken, errPublicKeyMarshal
	}

	conn.WritePacket((&encryptionRequestPacket{
		ServerID:    "",
		PublicKey:   asn1PublicKey,
		VerifyToken: verifyToken,
	}).Marshal())

	return verifyToken, nil
}

func handleEncryptionResponse(conn net.Conn, verifyToken []byte) ([]byte, error) {
	// C→S: Encryption Response

	// (Client auth via mojang web)
	// The client will generate a random 16-byte shared secret, to be used
	// with the AES/CFB8 stream cipher. It then encrypts it with the server's
	// public key (PKCS#1 v1.5 padded), and also encrypts the verify token
	// received. Both byte arrays in the Encryption Response packet will be
	// 128 bytes long because of the padding.

	received, err := conn.ReadPacket()
	if err != nil {
		return nil, errConnClosed
	}
	if received.ID != EncryptionResponse {
		return nil, errNoEncryptionResponse
	}

	response := encryptionResponsePacket{}
	if err = received.Scan(&response); err != nil {
		return nil, err
	}

	if int(response.VerifyTokenLength) != verifyTokenLength {
		return nil, errWrongTokenLength
	}

	// The server decrypts the shared secret and token using its private key,
	// and checks if the token is the same.
	secret, err := rsa.DecryptPKCS1v15(rand.Reader, global.PrivateKey, response.Secret)
	if err != nil {
		return nil, rsa.ErrDecryption
	}
	token, err := rsa.DecryptPKCS1v15(rand.Reader, global.PrivateKey, response.VerifyToken)
	if err != nil {
		return nil, rsa.ErrDecryption
	}
	if bytes.Compare(token, verifyToken) != 0 {
		return nil, errVerifyToken
	}

	return secret, nil
}

func handleLoginSuccess(conn net.Conn, username packet.String, uuid packet.String) error {
	conn.WritePacket((&loginSuccessPacket{
		UUID:     uuid,
		Username: username,
	}).Marshal())
	return nil
}
