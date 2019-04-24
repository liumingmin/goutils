package utils

import (
	"bytes"
	"encoding/binary"
	"io"

	"errors"
)

var (
	RPC_PACKET_FLAG = [2]byte{0xEF, 0xEE}
)

type PacketHeader struct {
	flag       [2]byte
	length     uint16
	ProtocolId uint16
}

type DataPacket struct {
	PacketHeader
	Data []byte
}

func (p *DataPacket) Unpack(reader io.Reader) error {
	var headerBytes [6]byte
	_, err := io.ReadAtLeast(reader, headerBytes[:], 6)
	if err != nil {
		return err
	}

	flagBytes := headerBytes[:2]
	if flagBytes[0] != RPC_PACKET_FLAG[0] || flagBytes[1] != RPC_PACKET_FLAG[1] {
		return errors.New("packet flag error")
	}

	lengthByte := headerBytes[2:4]
	protocolIdByte := headerBytes[4:6]

	binary.Read(bytes.NewBuffer(lengthByte[:]), binary.LittleEndian, &p.length)
	binary.Read(bytes.NewBuffer(protocolIdByte[:]), binary.LittleEndian, &p.ProtocolId)

	p.Data = make([]byte, p.length)
	reader.Read(p.Data)
	return nil
}

func (p *DataPacket) Pack(writer io.Writer) error {
	var err error

	var lengthByte [2]byte
	var protocolIdByte [2]byte

	p.length = uint16(len(p.Data))

	binary.LittleEndian.PutUint16(lengthByte[:], p.length)
	binary.LittleEndian.PutUint16(protocolIdByte[:], p.ProtocolId)

	flagBytes := []byte{RPC_PACKET_FLAG[0], RPC_PACKET_FLAG[1]}
	dataPack := append(flagBytes, lengthByte[:]...)
	dataPack = append(dataPack, protocolIdByte[:]...)
	dataPack = append(dataPack, p.Data...)

	err = binary.Write(writer, binary.BigEndian, &dataPack)
	return err
}
