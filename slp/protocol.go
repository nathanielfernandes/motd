package slp

import (
	"errors"
	"net"
)

const STATUS = 1
const LOGIN = 2

const SEGMENT_BITS = 0x7F
const CONTINUE_BIT = 0x80

type VarInt = int32
type VarLong = int64

func consumePacketId(conn net.Conn) error {
	// read and discard the first 2 bytes
	var firstTwoBytes = make([]byte, 2)
	_, err := conn.Read(firstTwoBytes)
	return err
}

func ReadStatusRequest(conn net.Conn) error {
	return consumePacketId(conn)
}

// read a varint from the connection
func readVarInt(conn net.Conn) (VarInt, error) {
	var value int32 = 0
	var position int = 0
	currentByte := make([]byte, 1)

	for {

		if _, err := conn.Read(currentByte); err != nil {
			return 0, err
		}

		value |= int32(currentByte[0]&SEGMENT_BITS) << position

		if currentByte[0]&CONTINUE_BIT == 0 {
			break
		}

		position += 7

		if position > 32 {
			return 0, errors.New("VarInt too big")
		}
	}

	return value, nil
}

// read varlong from the connection
func readVarLong(conn net.Conn) (VarLong, error) {
	var value int64 = 0
	var position int = 0
	var currentByte byte

	for {

		if _, err := conn.Read([]byte{currentByte}); err != nil {
			return 0, err
		}

		value |= int64(currentByte&SEGMENT_BITS) << position

		if currentByte&CONTINUE_BIT == 0 {
			break
		}

		position += 7

		if position > 64 {
			return 0, errors.New("VarLong too big")
		}
	}

	return value, nil
}

// read an unsigned short from the connection
func readUnsignedShort(conn net.Conn) (uint16, error) {
	var shortBytes = make([]byte, 2)

	if _, err := conn.Read(shortBytes); err != nil {
		return 0, err
	}

	return uint16(shortBytes[0])<<8 | uint16(shortBytes[1]), nil
}

// read a string from the connection
func readString(conn net.Conn) (string, error) {
	var length int32
	var err error

	length, err = readVarInt(conn)
	if err != nil {
		return "", err
	}

	var stringBytes = make([]byte, length)

	if _, err = conn.Read(stringBytes); err != nil {
		return "", err
	}

	return string(stringBytes), nil
}

// Create a varint from an int32
func NewVarInt(value int32) []byte {
	var bytes = make([]byte, 0)

	for {
		var temp = byte(value & SEGMENT_BITS)

		value >>= 7
		if value != 0 {
			temp |= CONTINUE_BIT
		}

		bytes = append(bytes, temp)

		if value == 0 {
			break
		}
	}

	return bytes
}

func NewVarLong(value int64) []byte {
	var bytes = make([]byte, 0)

	for {
		var temp = byte(value & SEGMENT_BITS)

		value >>= 7
		if value != 0 {
			temp |= CONTINUE_BIT
		}

		bytes = append(bytes, temp)

		if value == 0 {
			break
		}
	}

	return bytes
}

func NewPacketId(id int32) []byte {
	return NewVarInt(id)
}

func NewBool(value bool) []byte {
	if value {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func readBool(conn net.Conn) (bool, error) {
	var boolByte byte

	if _, err := conn.Read([]byte{boolByte}); err != nil {
		return false, err
	}

	return boolByte != 0, nil
}
