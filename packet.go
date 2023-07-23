package openblive

import (
	"bytes"
	"encoding/binary"
)

const (
	// OpHandshake handshake
	OpHandshake = 0
	// OpHandshakeReply handshake reply
	OpHandshakeReply = 1

	// OpHeartbeat heartbeat
	OpHeartbeat = 2
	// OpHeartbeatReply heartbeat reply
	OpHeartbeatReply = 3

	// OpSendMsg send message.
	OpSendMsg = 4
	// OpSendMsgReply  send message reply
	OpSendMsgReply = 5

	// OpDisconnectReply disconnect reply
	OpDisconnectReply = 6

	// OpAuth auth connnect
	OpAuth = 7
	// OpAuthReply auth connect reply
	OpAuthReply = 8

	// OpRaw  raw message
	OpRaw = 9

	// OpProtoReady proto ready
	OpProtoReady = 10
	// OpProtoFinish proto finish
	OpProtoFinish = 11

	// OpChangeRoom change room
	OpChangeRoom = 12
	// OpChangeRoomReply change room reply
	OpChangeRoomReply = 13

	// OpRegister register operation
	OpRegister = 14
	// OpRegisterReply register operation
	OpRegisterReply = 15

	// OpUnregister unregister operation
	OpUnregister = 16
	// OpUnregisterReply unregister operation reply
	OpUnregisterReply = 17

	// MinBusinessOp min business operation
	MinBusinessOp = 1000
	// MaxBusinessOp max business operation
	MaxBusinessOp = 10000
)

type WsHeader struct {
	PacketLength    uint32
	HeaderLength    uint16
	ProtocolVersion uint16
	Operation       uint32
	Sequence        uint32
}

type WsPacket struct {
	Header WsHeader
	Data   []byte
}

func ResolveWSPacket(data []byte) (WsPacket, bool) {
	if len(data) < 16 {
		return WsPacket{}, false
	}
	header := WsHeader{
		PacketLength:    binary.BigEndian.Uint32(data[0:4]),
		HeaderLength:    binary.BigEndian.Uint16(data[4:6]),
		ProtocolVersion: binary.BigEndian.Uint16(data[6:8]),
		Operation:       binary.BigEndian.Uint32(data[8:12]),
		Sequence:        binary.BigEndian.Uint32(data[12:16]),
	}
	return WsPacket{
		Header: header,
		Data:   data[header.HeaderLength:header.PacketLength],
	}, true
}

func MakeWSPacket(operation int, data []byte) []byte {
	headerBytes := new(bytes.Buffer)
	header := []interface{}{
		uint32(len(data) + 16),
		uint16(16),
		uint16(0),
		uint32(operation),
		uint32(1),
	}
	for _, v := range header {
		err := binary.Write(headerBytes, binary.BigEndian, v)
		if err != nil {
		}
	}
	return append(headerBytes.Bytes(), data...)
}
