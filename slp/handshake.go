package slp

import (
	"fmt"
	"net"
)

type Handshake struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       int32
}

func (h *Handshake) ReadFromConn(conn net.Conn) error {
	err := consumePacketId(conn)
	if err != nil {
		return err
	}

	if h.ProtocolVersion, err = readVarInt(conn); err != nil {
		return err
	}

	if h.ServerAddress, err = readString(conn); err != nil {
		return err
	}

	if h.ServerPort, err = readUnsignedShort(conn); err != nil {
		return err
	}

	if h.NextState, err = readVarInt(conn); err != nil {
		return err
	}

	return nil
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
