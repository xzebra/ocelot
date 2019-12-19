package packet

// Packet uncompressed format
type Packet struct {
	// Length of packet data + length of the packet ID
	Length VarInt
	// ID of the packet
	ID VarInt
	// Data depends on the connection state and packet ID
	Data ByteArray
}

// CompressedPacket is used when compression is set
type CompressedPacket struct {
	// Length is the length of Data Length + compressed
	// length of (Packet ID + Data)
	Length VarInt
	// DataLength is the length of uncompressed(ID + Data) or 0
	// If DataLength is set to 0, the packet is uncompressed
	DataLength VarInt
	// zlib compressed packet ID
	ID VarInt
	// zlib compressed packet data
	Data ByteArray
}
