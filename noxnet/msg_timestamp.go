package noxnet

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgTimestamp{}, false)
	netmsg.Register(&MsgFullTimestamp{}, false)
	netmsg.Register(&MsgRateChange{}, false)
}

type Timestamp uint32

func (t Timestamp) Dur(fps int) time.Duration {
	return time.Duration(t) * time.Second / time.Duration(fps)
}

func (t *Timestamp) Set(v Timestamp) {
	if v > *t {
		*t = v
	}
}

func (t *Timestamp) Set16(v uint16) {
	ts := uint32(*t)
	ts16 := uint16(ts & 0xFFFF)
	overflow := (ts16 >= 0xC000) && (v < 0x4000)
	if !overflow && v < ts16 {
		return // out of order
	}
	ts = (ts & 0xFFFF0000) | uint32(v)
	if overflow {
		ts += 0x10000
	}
	*t = Timestamp(ts)
}

type MsgTimestamp struct {
	T uint16
}

func (*MsgTimestamp) NetOp() netmsg.Op {
	return netmsg.MSG_TIMESTAMP
}

func (*MsgTimestamp) EncodeSize() int {
	return 2
}

func (m *MsgTimestamp) Encode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[:2], m.T)
	return 2, nil
}

func (m *MsgTimestamp) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	m.T = binary.LittleEndian.Uint16(data[:2])
	return 2, nil
}

type MsgFullTimestamp struct {
	T Timestamp
}

func (*MsgFullTimestamp) NetOp() netmsg.Op {
	return netmsg.MSG_FULL_TIMESTAMP
}

func (*MsgFullTimestamp) EncodeSize() int {
	return 4
}

func (m *MsgFullTimestamp) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[:4], uint32(m.T))
	return 4, nil
}

func (m *MsgFullTimestamp) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.T = Timestamp(binary.LittleEndian.Uint32(data[:4]))
	return 4, nil
}

type MsgRateChange struct {
	Rate byte
}

func (*MsgRateChange) NetOp() netmsg.Op {
	return netmsg.MSG_RATE_CHANGE
}

func (*MsgRateChange) EncodeSize() int {
	return 1
}

func (p *MsgRateChange) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Rate
	return 1, nil
}

func (p *MsgRateChange) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Rate = data[0]
	return 1, nil
}
