package mapsend

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgMapSendStart{}, false)
	netmsg.Register(&MsgMapSendAbort{}, false)
	netmsg.Register(&MsgMapSendPacket{}, true)
	netmsg.Register(&MsgMapReceived{}, false)
}

type MsgMapSendStart struct {
	Unk1    [3]byte
	MapSize uint32
	MapName binenc.String
}

func (*MsgMapSendStart) NetOp() netmsg.Op {
	return netmsg.MSG_MAP_SEND_START
}

func (*MsgMapSendStart) EncodeSize() int {
	return 87
}

func (p *MsgMapSendStart) Encode(data []byte) (int, error) {
	if len(data) < 87 {
		return 0, io.ErrShortBuffer
	}
	copy(data[0:3], p.Unk1[:])
	binary.LittleEndian.PutUint32(data[3:7], p.MapSize)
	p.MapName.Encode(data[7:87])
	return 87, nil
}

func (p *MsgMapSendStart) Decode(data []byte) (int, error) {
	if len(data) < 87 {
		return 0, io.ErrUnexpectedEOF
	}
	copy(p.Unk1[:], data[0:3])
	p.MapSize = binary.LittleEndian.Uint32(data[3:7])
	p.MapName.Decode(data[7:87])
	return 87, nil
}

type MsgMapSendAbort struct {
	Code byte
}

func (*MsgMapSendAbort) NetOp() netmsg.Op {
	return netmsg.MSG_MAP_SEND_ABORT
}

func (*MsgMapSendAbort) EncodeSize() int {
	return 1
}

func (p *MsgMapSendAbort) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Code
	return 1, nil
}

func (p *MsgMapSendAbort) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Code = data[0]
	return 1, nil
}

type MsgMapSendPacket struct {
	Unk   byte
	Block uint16
	Data  []byte
}

func (*MsgMapSendPacket) NetOp() netmsg.Op {
	return netmsg.MSG_MAP_SEND_PACKET
}

func (p *MsgMapSendPacket) EncodeSize() int {
	return 5 + len(p.Data)
}

func (p *MsgMapSendPacket) Encode(data []byte) (int, error) {
	if len(data) < 5+len(p.Data) {
		return 0, io.ErrShortBuffer
	}
	if len(p.Data) > math.MaxUint16 {
		return 0, errors.New("map packet too large")
	}
	data[0] = p.Unk
	binary.LittleEndian.PutUint16(data[1:3], p.Block)
	binary.LittleEndian.PutUint16(data[3:5], uint16(len(p.Data)))
	copy(data[5:5+len(p.Data)], p.Data)
	return 5 + len(p.Data), nil
}

func (p *MsgMapSendPacket) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Unk = data[0]
	p.Block = binary.LittleEndian.Uint16(data[1:3])
	sz := int(binary.LittleEndian.Uint16(data[3:5]))
	data = data[5:]
	if len(data) < sz {
		return 0, io.ErrUnexpectedEOF
	}
	p.Data = make([]byte, sz)
	copy(p.Data, data)
	return 5 + sz, nil
}

type MsgMapReceived struct {
}

func (*MsgMapReceived) NetOp() netmsg.Op {
	return netmsg.MSG_RECEIVED_MAP
}

func (*MsgMapReceived) EncodeSize() int {
	return 1
}

func (*MsgMapReceived) Encode(data []byte) (int, error) {
	return 0, nil
}

func (*MsgMapReceived) Decode(data []byte) (int, error) {
	return 0, nil
}
