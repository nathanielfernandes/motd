package slp

import (
	"fmt"
	"io"
	"net"
)

type Handshake struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       int32
}

func (h *Handshake) Read(r io.Reader) (err error) {
	if err = ReadStatusRequest(r); err != nil {
		return
	}

	if h.ProtocolVersion, err = ReadVarInt(r); err != nil {
		return
	}

	if h.ServerAddress, err = ReadString(r); err != nil {
		return
	}

	if h.ServerPort, err = ReadUnsignedShort(r); err != nil {
		return
	}

	if h.NextState, err = ReadVarInt(r); err != nil {
		return
	}

	return
}

func (h Handshake) Write(w io.Writer) (err error) {
	packet := NewPacket(0x00)

	if _, err = WriteVarInt(&packet, h.ProtocolVersion); err != nil {
		return
	}

	if _, err = WriteString(&packet, h.ServerAddress); err != nil {
		return
	}

	if _, err = WriteUnsignedShort(&packet, h.ServerPort); err != nil {
		return
	}

	if _, err = WriteVarInt(&packet, h.NextState); err != nil {
		return
	}

	_, err = packet.WritePacket(w)

	return
}

func (h Handshake) PrettyPrint() {
	var state string
	switch h.NextState {
	case 1:
		state = "Status"
	case 2:
		state = "Login"
	default:
		state = "Unknown"
	}

	fmt.Printf("Handshake Established:\n  Protocol Version: %d\n  Server Address: %s\n  Server Port: %d\n  Next State: %d (%s)\n", h.ProtocolVersion, h.ServerAddress, h.ServerPort, h.NextState, state)
}

func ListenAndPong(conn net.Conn) (err error) {
	if err = ReadStatusRequest(conn); err != nil {
		return
	}

	pingId, err := ReadLong(conn)
	if err != nil {
		return
	}

	_, err = WriteLong(conn, pingId)

	return
}

func PingAndClose(conn net.Conn) (err error) {
	packet := NewPacket(0x01)

	if _, err = WriteLong(&packet, 1234); err != nil {
		return
	}

	if _, err = packet.WritePacket(conn); err != nil {
		return
	}

	if err = ReadStatusRequest(conn); err != nil {
		return
	}

	// silly tbh
	ListenAndPong(conn)

	return conn.Close()
}
