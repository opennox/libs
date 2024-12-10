package noxnet

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgClientPing{}, false)
	netmsg.Register(&MsgClientPong{}, false)
	netmsg.Register(&MsgServerPing{}, false)
	netmsg.Register(&MsgServerPong{}, false)
	netmsg.Register(&MsgSpeed{}, false)
	netmsg.Register(&MsgKeepAlive{}, false)
}

type MsgClientPing struct {
	T time.Duration // in ms
}

func (m *MsgClientPing) Pong() *MsgClientPong {
	return &MsgClientPong{
		T: m.T,
	}
}

func (*MsgClientPing) NetOp() netmsg.Op {
	return netmsg.MSG_CLIENT_PING
}

func (*MsgClientPing) EncodeSize() int {
	return 4
}

func (m *MsgClientPing) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:], uint32(m.T/time.Millisecond))
	return 4, nil
}

func (m *MsgClientPing) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.T = time.Duration(binary.LittleEndian.Uint32(data[0:])) * time.Millisecond
	return 4, nil
}

type MsgClientPong struct {
	T time.Duration // in ms
}

func (*MsgClientPong) NetOp() netmsg.Op {
	return netmsg.MSG_CLIENT_PONG
}

func (*MsgClientPong) EncodeSize() int {
	return 4
}

func (m *MsgClientPong) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:], uint32(m.T/time.Millisecond))
	return 4, nil
}

func (m *MsgClientPong) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.T = time.Duration(binary.LittleEndian.Uint32(data[0:])) * time.Millisecond
	return 4, nil
}

type MsgServerPing struct {
	Ind byte
	T   time.Duration // in ms
}

func (m *MsgServerPing) Pong() *MsgServerPong {
	return &MsgServerPong{
		Ind: m.Ind,
		T:   m.T,
	}
}

func (*MsgServerPing) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_PING
}

func (*MsgServerPing) EncodeSize() int {
	return 5
}

func (m *MsgServerPing) Encode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, io.ErrShortBuffer
	}
	data[0] = m.Ind
	binary.LittleEndian.PutUint32(data[1:], uint32(m.T/time.Millisecond))
	return 5, nil
}

func (m *MsgServerPing) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Ind = data[0]
	m.T = time.Duration(binary.LittleEndian.Uint32(data[1:])) * time.Millisecond
	return 5, nil
}

type MsgServerPong struct {
	Ind byte
	T   time.Duration // in ms
}

func (*MsgServerPong) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_PONG
}

func (*MsgServerPong) EncodeSize() int {
	return 5
}

func (m *MsgServerPong) Encode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, io.ErrShortBuffer
	}
	data[0] = m.Ind
	binary.LittleEndian.PutUint32(data[1:], uint32(m.T/time.Millisecond))
	return 5, nil
}

func (m *MsgServerPong) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Ind = data[0]
	m.T = time.Duration(binary.LittleEndian.Uint32(data[1:])) * time.Millisecond
	return 5, nil
}

type MsgSpeed struct {
	Speed int32 // 256 KB / T ms (or -1)
}

func (*MsgSpeed) NetOp() netmsg.Op {
	return netmsg.MSG_SPEED
}

func (*MsgSpeed) EncodeSize() int {
	return 4
}

func (m *MsgSpeed) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:], uint32(m.Speed))
	return 4, nil
}

func (m *MsgSpeed) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Speed = int32(binary.LittleEndian.Uint32(data[0:]))
	return 4, nil
}

type MsgKeepAlive struct {
}

func (*MsgKeepAlive) NetOp() netmsg.Op {
	return netmsg.MSG_KEEP_ALIVE
}

func (*MsgKeepAlive) EncodeSize() int {
	return 0
}

func (m *MsgKeepAlive) Encode(data []byte) (int, error) {
	return 0, nil
}

func (m *MsgKeepAlive) Decode(data []byte) (int, error) {
	return 0, nil
}
