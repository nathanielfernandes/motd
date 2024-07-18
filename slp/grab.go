package slp

import (
	"net"
)

// get a status from a server
func GrabStatus(address string) (r StatusResponse, err error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	// write handshake
	hs := Handshake{
		ProtocolVersion: 767,
		ServerAddress:   "",
		ServerPort:      25565,
		NextState:       NEXTSTATE_STATUS,
	}

	if err = hs.Write(conn); err != nil {
		return
	}

	// write status request
	if _, err = NewPacket(0x00).WritePacket(conn); err != nil {
		return
	}

	// read status response
	if err = r.Read(conn); err != nil {
		panic(err)
	}

	go PingAndClose(conn)

	return
}
