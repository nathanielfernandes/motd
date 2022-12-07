package slp

import (
	"encoding/json"
	"net"
)

type LoginStart struct {
	Name            string
	HasSigData      bool
	Timestamp       VarLong
	PublicKeyLength VarInt
	PublicKey       []byte
	SignatureLength VarInt
	Signature       []byte
	HasPlayerUUID   bool
	PlayerUUID      string
}

func (ls *LoginStart) ReadFromConn(conn net.Conn) error {
	err := consumePacketId(conn)
	if err != nil {
		return err
	}

	if ls.Name, err = readString(conn); err != nil {
		return err
	}

	if ls.HasSigData, err = readBool(conn); err != nil {
		return err
	}

	if ls.HasSigData {
		if ls.Timestamp, err = readVarLong(conn); err != nil {
			return err
		}

		if ls.PublicKeyLength, err = readVarInt(conn); err != nil {
			return err
		}

		ls.PublicKey = make([]byte, ls.PublicKeyLength)
		if _, err = conn.Read(ls.PublicKey); err != nil {
			return err
		}

		if ls.SignatureLength, err = readVarInt(conn); err != nil {
			return err
		}

		ls.Signature = make([]byte, ls.SignatureLength)
		if _, err = conn.Read(ls.Signature); err != nil {
			return err
		}

		if ls.HasPlayerUUID, err = readBool(conn); err != nil {
			return err
		}

		if ls.HasPlayerUUID {
			if ls.PlayerUUID, err = readString(conn); err != nil {
				return err
			}
		}
	}

	return nil
}

type Disconnect struct {
	Reason Chat
}

func (d *Disconnect) WriteToConn(conn net.Conn) error {
	packet := NewPacket(0x00)

	chatJson, err := json.Marshal(d.Reason)
	if err != nil {
		return err
	}

	packet.WriteJsonBytes(chatJson)

	return packet.WriteToConn(conn)
}

func DisconnectWithStringMsg(msg string) Disconnect {
	return Disconnect{Reason: Chat{Text: msg}}
}

func DisconnectWithChatMsg(chat Chat) Disconnect {
	return Disconnect{Reason: chat}
}
