package slp

import (
	"io"
)

type LoginStart struct {
	Name string
	UUID [16]byte
}

func (ls *LoginStart) Read(r io.Reader) (err error) {
	if err = ReadStatusRequest(r); err != nil {
		return
	}

	if ls.Name, err = ReadString(r); err != nil {
		return
	}

	if _, err = r.Read(ls.UUID[:]); err != nil {
		return
	}

	return
}

func (ls *LoginStart) WriteSuccess(w io.Writer) (err error) {
	packet := NewPacket(0x02)

	if _, err = WriteUUID(&packet, ls.UUID); err != nil {
		return
	}

	if _, err = WriteString(&packet, ls.Name); err != nil {
		return
	}

	// num properties
	if _, err = WriteVarInt(&packet, 0); err != nil {
		return
	}

	if _, err = WriteBool(&packet, false); err != nil {
		return
	}

	_, err = packet.WritePacket(w)

	return
}
