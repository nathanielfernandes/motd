package slp

import "net"

type Packet []byte

func (p *Packet) Write(data []byte) {
	*p = append(*p, data...)
}

func NewPacket(packetID byte) Packet {
	return Packet{packetID}
}

func (p *Packet) WriteVarInt(value VarInt) {
	p.Write(NewVarInt(value))
}

func (p *Packet) WriteVarLong(value VarLong) {
	p.Write(NewVarLong(value))
}

func (p *Packet) WriteString(value string) {
	p.Write([]byte(value))
}

func (p *Packet) WriteBool(value bool) {
	p.Write(NewBool(value))
}

func (p *Packet) WriteJsonBytes(data []byte) {
	p.WriteVarInt(VarInt(len(data)))
	p.Write(data)
}

func (p *Packet) WriteToConn(conn net.Conn) error {
	if _, err := conn.Write(NewVarInt(VarInt(len(*p)))); err != nil {
		return err
	}
	_, err := conn.Write(*p)
	return err
}
