package slp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const SEGMENT_BITS = 0x7F
const CONTINUE_BIT = 0x80

func ReadStatusRequest(r io.Reader) (err error) {
	var length int32
	if length, err = ReadVarInt(r); err != nil {
		return errors.New("Error reading status request: " + err.Error())
	}
	fmt.Println("Length:", length)

	var packetId int32
	if packetId, err = ReadVarInt(r); err != nil {
		return errors.New("Error reading status request: " + err.Error())
	}
	fmt.Println("Packet ID:", packetId)

	return nil
}

func ReadVarInt(r io.Reader) (int32, error) {
	var value int32 = 0
	var position int = 0
	currentByte := make([]byte, 1)

	for {
		if _, err := r.Read(currentByte); err != nil {
			return 0, errors.New("Error reading VarInt: " + err.Error())
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

func ReadLong(r io.Reader) (int64, error) {
	data := make([]byte, 8)
	if _, err := r.Read(data); err != nil {
		return 0, errors.New("Error reading long: " + err.Error())
	}
	return int64(data[0])<<56 | int64(data[1])<<48 | int64(data[2])<<40 | int64(data[3])<<32 | int64(data[4])<<24 | int64(data[5])<<16 | int64(data[6])<<8 | int64(data[7]), nil
}

func ReadVarLong(r io.Reader) (int64, error) {
	var value int64 = 0
	var position int = 0
	currentByte := make([]byte, 1)

	for {
		if _, err := r.Read(currentByte); err != nil {
			return 0, errors.New("Error reading VarLong: " + err.Error())
		}

		value |= int64(currentByte[0]&SEGMENT_BITS) << position

		if currentByte[0]&CONTINUE_BIT == 0 {
			break
		}

		position += 7

		if position > 64 {
			return 0, errors.New("VarLong too big")
		}
	}

	return value, nil
}

func ReadUnsignedShort(r io.Reader) (uint16, error) {
	data := make([]byte, 2)
	if _, err := r.Read(data); err != nil {
		return 0, errors.New("Error reading unsigned short: " + err.Error())
	}
	return uint16(data[0])<<8 | uint16(data[1]), nil
}

func ReadBool(r io.Reader) (bool, error) {
	data := make([]byte, 1)
	if _, err := r.Read(data); err != nil {
		return false, errors.New("Error reading bool: " + err.Error())
	}
	return data[0] == 1, nil
}

// keep checking the read bytes until we reach the length
func ReadByteSeq(r io.Reader) ([]byte, error) {
	length, err := ReadVarInt(r)
	if err != nil {
		return nil, errors.New("Error reading data length: " + err.Error())
	}
	var data []byte
	if length > 1024 {
		data = make([]byte, 0, length)

		reader := io.LimitReader(r, int64(length))
		buffer := make([]byte, 1024) // 1kb buffer

		for {
			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}

				return nil, errors.New("Error reading data: " + err.Error())
			}
			data = append(data, buffer[:n]...)
		}
	} else {
		data = make([]byte, length)
		if _, err = r.Read(data); err != nil {
			return nil, errors.New("Error reading data: " + err.Error())
		}
	}

	return data, nil
}

func ReadString(r io.Reader) (string, error) {
	data, err := ReadByteSeq(r)
	if err != nil {
		return "", errors.New("Error reading string: " + err.Error())
	}
	return string(data), nil
}

func ReadJsonBytes[T any](r io.Reader, target *T) error {
	data, err := ReadByteSeq(r)
	if err != nil {
		return errors.New("Error reading JSON: " + err.Error())
	}

	if err = json.Unmarshal(data, target); err != nil {
		return errors.New("Error unmarshalling JSON: " + err.Error())
	}

	return nil
}

func WriteVarInt(w io.Writer, value int32) (int, error) {
	var buffer []byte
	for {
		temp := byte(value & SEGMENT_BITS)
		value >>= 7

		if value != 0 {
			temp |= CONTINUE_BIT
		}

		buffer = append(buffer, temp)
		if value == 0 {
			break
		}
	}

	return w.Write(buffer)
}

func WriteVarLong(w io.Writer, value int64) (int, error) {
	var buffer []byte
	for {
		temp := byte(value & SEGMENT_BITS)
		value >>= 7

		if value != 0 {
			temp |= CONTINUE_BIT
		}

		buffer = append(buffer, temp)
		if value == 0 {
			break
		}
	}

	return w.Write(buffer)
}

func WriteUnsignedShort(w io.Writer, value uint16) (int, error) {
	return w.Write([]byte{byte(value >> 8), byte(value)})
}

func WritePacketId(w io.Writer, packetID int32) (int, error) {
	return w.Write([]byte{byte(packetID)})
}

func WriteBool(w io.Writer, value bool) (int, error) {
	if value {
		return w.Write([]byte{1})
	} else {
		return w.Write([]byte{0})
	}
}

func WriteLong(w io.Writer, value int64) (int, error) {
	return w.Write([]byte{byte(value >> 56), byte(value >> 48), byte(value >> 40), byte(value >> 32), byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

func WriteByteSeq(w io.Writer, data []byte) (n int, err error) {
	if n, err = WriteVarInt(w, int32(len(data))); err != nil {
		return
	}

	return w.Write(data)
}

func WriteUUID(w io.Writer, uuid [16]byte) (int, error) {
	return w.Write(uuid[:])
}

func WriteString(w io.Writer, value string) (int, error) {
	return WriteByteSeq(w, []byte(value))
}

func WriteJsonBytes(w io.Writer, data []byte) (int, error) {
	return WriteByteSeq(w, data)
}
