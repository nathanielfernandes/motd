package slp

import (
	"encoding/json"
	"net"
)

type StatusResponse struct {
	Version            Version     `json:"version"`
	Players            Players     `json:"players"`
	Description        Description `json:"description"`
	Favicon            string      `json:"favicon"`
	PreviewsChat       bool        `json:"previewsChat"`
	EnforcesSecureChat bool        `json:"enforcesSecureChat"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

func SpoofVersion(text string) Version {
	return Version{
		Name:     text,
		Protocol: 735, // old version so that the name is shown
	}
}

type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []Sample `json:"sample"`
}

type Sample struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Description struct {
	Text string `json:"text"`
}

func (sr *StatusResponse) WriteToConn(conn net.Conn) error {
	packet := NewPacket(0x00)

	jsonResponseBytes, err := json.Marshal(sr)
	if err != nil {
		return err
	}

	packet.WriteJsonBytes(jsonResponseBytes)

	return packet.WriteToConn(conn)
}

func PingPong(conn net.Conn) error {
	// read ping
	consumePacketId(conn)
	pingId, err := readVarLong(conn)
	if err != nil {
		return err
	}

	// write pong
	_, err = conn.Write(NewPacketId(0x01))
	if err != nil {
		return err
	}
	_, err = conn.Write(NewVarLong(pingId))

	return err
}
