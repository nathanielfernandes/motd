package slp

import "io"

// helper type for writing packets
type Packet []byte

func NewPacket(packetID byte) Packet {
	return Packet{packetID}
}

func (p *Packet) Write(data []byte) (int, error) {
	*p = append(*p, data...)
	return len(data), nil
}

func (p Packet) WritePacket(w io.Writer) (int, error) {
	return WriteByteSeq(w, p)
}
