package utils

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ControlPacket struct {
	Length     uint16
	ProtocolId uint16
	Data       []byte
}

func (p *ControlPacket) Unpack(reader io.Reader) {
	var lengthByte [2]byte
	var protocolIdByte [2]byte

	reader.Read(lengthByte[:])
	reader.Read(protocolIdByte[:])

	binary.Read(bytes.NewBuffer(lengthByte[:]), binary.LittleEndian, &p.Length)
	binary.Read(bytes.NewBuffer(protocolIdByte[:]), binary.LittleEndian, &p.ProtocolId)

	p.Data = make([]byte, p.Length-4)
	reader.Read(p.Data)
}

func (p *ControlPacket) Pack(writer io.Writer) error {
	var err error

	var lengthByte [2]byte
	var protocolIdByte [2]byte

	p.Length = uint16(len(p.Data) + 4)

	binary.LittleEndian.PutUint16(lengthByte[:], p.Length)
	binary.LittleEndian.PutUint16(protocolIdByte[:], p.ProtocolId)

	dataPack := append(protocolIdByte[:], p.Data...)
	dataPack = append(lengthByte[:], dataPack...)

	err = binary.Write(writer, binary.BigEndian, &dataPack)
	return err
}

func ReadPacketLength16(r io.Reader) (uint16, error) {
	var lengthByte [2]byte
	_, err := io.ReadAtLeast(r, lengthByte[:], 2)
	if err != nil {
		return 0, err
	}

	var length uint16
	binary.Read(bytes.NewBuffer(lengthByte[:]), binary.LittleEndian, &length)

	return length, nil
}

func ReadPacketLength(r io.Reader) (uint32, error) {
	var lengthByte [4]byte
	_, err := io.ReadAtLeast(r, lengthByte[:], 4)
	if err != nil {
		return 0, err
	}

	var length uint32
	binary.Read(bytes.NewBuffer(lengthByte[:]), binary.LittleEndian, &length)

	return length, nil
}
