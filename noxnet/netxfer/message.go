package netxfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/opennox/libs/binenc"
)

type Op byte

const (
	OpStart  = Op(0)
	OpAccept = Op(1)
	OpData   = Op(2)
	OpAck    = Op(3)
	OpDone   = Op(4)
	OpCancel = Op(5)
	OpAbort  = Op(6)
)

type Action byte
type SendID byte
type RecvID byte
type Chunk uint16
type Error byte

func (e Error) Error() string {
	if int(e) < len(errTexts) {
		if s := errTexts[e]; s != "" {
			return "netxfer: " + s
		}
	}
	return fmt.Sprintf("netxfer: unknown error (%d)", e)
}

const (
	ErrClosed      = Error(1)
	ErrSendTimeout = Error(2)
	ErrRecvTimeout = Error(3)
)

var errTexts = []string{
	ErrClosed:      "send cancelled",
	ErrSendTimeout: "send timeout",
	ErrRecvTimeout: "recv timeout",
}

type Msg interface {
	XferOp() Op
	EncodeSize() int
	Encode(data []byte) (int, error)
	Decode(data []byte) (int, error)
}

type MsgStart struct {
	Act    Action
	Unk1   byte
	Size   uint32
	Type   binenc.String
	SendID SendID
	Unk5   [3]byte
}

func (*MsgStart) XferOp() Op {
	return OpStart
}

func (*MsgStart) EncodeSize() int {
	return 138
}

func (m *MsgStart) Encode(data []byte) (int, error) {
	if len(data) < 138 {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(m.Act)
	data[1] = m.Unk1
	binary.LittleEndian.PutUint32(data[2:6], m.Size)
	m.Type.Encode(data[6:134])
	data[134] = byte(m.SendID)
	copy(data[135:138], m.Unk5[:])
	return 138, nil
}

func (m *MsgStart) Decode(data []byte) (int, error) {
	if len(data) < 138 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Act = Action(data[0])
	m.Unk1 = data[1]
	m.Size = binary.LittleEndian.Uint32(data[2:6])
	m.Type.Decode(data[6:134])
	m.SendID = SendID(data[134])
	copy(m.Unk5[:], data[135:138])
	return 138, nil
}

type MsgAccept struct {
	RecvID RecvID // use this to send data
	SendID SendID // accepted send stream
}

func (*MsgAccept) XferOp() Op {
	return OpAccept
}

func (*MsgAccept) EncodeSize() int {
	return 2
}

func (m *MsgAccept) Encode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(m.RecvID)
	data[1] = byte(m.SendID)
	return 2, nil
}

func (m *MsgAccept) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	m.RecvID = RecvID(data[0])
	m.SendID = SendID(data[1])
	return 2, nil
}

type MsgData struct {
	RecvID RecvID
	Token  byte
	Chunk  Chunk
	Data   []byte
}

func (*MsgData) XferOp() Op {
	return OpData
}

func (m *MsgData) EncodeSize() int {
	return 6 + len(m.Data)
}

func (m *MsgData) Encode(data []byte) (int, error) {
	if len(data) < 6+len(m.Data) {
		return 0, io.ErrShortBuffer
	}
	if len(m.Data) > math.MaxUint16 {
		return 0, errors.New("xfer packet too large")
	}
	data[0] = byte(m.RecvID)
	data[1] = m.Token
	binary.LittleEndian.PutUint16(data[2:4], uint16(m.Chunk))
	binary.LittleEndian.PutUint16(data[4:6], uint16(len(m.Data)))
	copy(data[6:], m.Data)
	return 6 + len(m.Data), nil
}

func (m *MsgData) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, io.ErrUnexpectedEOF
	}
	m.RecvID = RecvID(data[0])
	m.Token = data[1]
	m.Chunk = Chunk(binary.LittleEndian.Uint16(data[2:4]))
	sz := int(binary.LittleEndian.Uint16(data[4:6]))
	data = data[6:]
	if sz > len(data) {
		return 0, io.ErrUnexpectedEOF
	}
	m.Data = make([]byte, sz)
	copy(m.Data, data[:sz])
	return 6 + sz, nil
}

type MsgAck struct {
	RecvID RecvID
	Token  byte
	Chunk  Chunk
}

func (*MsgAck) XferOp() Op {
	return OpAck
}

func (*MsgAck) EncodeSize() int {
	return 4
}

func (m *MsgAck) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(m.RecvID)
	data[1] = m.Token
	binary.LittleEndian.PutUint16(data[2:4], uint16(m.Chunk))
	return 4, nil
}

func (m *MsgAck) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.RecvID = RecvID(data[0])
	m.Token = data[1]
	m.Chunk = Chunk(binary.LittleEndian.Uint16(data[2:4]))
	return 4, nil
}

type MsgDone struct {
	RecvID RecvID
}

func (*MsgDone) XferOp() Op {
	return OpDone
}

func (*MsgDone) EncodeSize() int {
	return 1
}

func (m *MsgDone) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(m.RecvID)
	return 1, nil
}

func (m *MsgDone) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	m.RecvID = RecvID(data[0])
	return 1, nil
}

type MsgCancel struct {
	RecvID RecvID
	Reason Error
}

func (*MsgCancel) XferOp() Op {
	return OpCancel
}

func (*MsgCancel) EncodeSize() int {
	return 2
}

func (m *MsgCancel) Encode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(m.RecvID)
	data[1] = byte(m.Reason)
	return 2, nil
}

func (m *MsgCancel) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	m.RecvID = RecvID(data[0])
	m.Reason = Error(data[1])
	return 2, nil
}

type MsgAbort struct {
	RecvID RecvID
	Reason Error
}

func (*MsgAbort) XferOp() Op {
	return OpAbort
}

func (*MsgAbort) EncodeSize() int {
	return 2
}

func (m *MsgAbort) Encode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(m.RecvID)
	data[1] = byte(m.Reason)
	return 2, nil
}

func (m *MsgAbort) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	m.RecvID = RecvID(data[0])
	m.Reason = Error(data[1])
	return 2, nil
}

func EncodeSize(m Msg) int {
	return 1 + m.EncodeSize()
}

func Encode(data []byte, m Msg) (int, error) {
	if len(data) < 1+m.EncodeSize() {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(m.XferOp())
	n, err := m.Encode(data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

func Decode(data []byte) (Msg, int, error) {
	if len(data) < 1 {
		return nil, 0, io.ErrUnexpectedEOF
	}
	typ := Op(data[0])
	var m Msg
	switch typ {
	case OpStart:
		m = &MsgStart{}
	case OpAccept:
		m = &MsgAccept{}
	case OpCancel:
		m = &MsgCancel{}
	case OpAbort:
		m = &MsgAbort{}
	case OpData:
		m = &MsgData{}
	case OpAck:
		m = &MsgAck{}
	case OpDone:
		m = &MsgDone{}
	default:
		return nil, 0, fmt.Errorf("unexpected xfer message subtype: 0x%x", typ)
	}
	n, err := m.Decode(data[1:])
	if err != nil {
		return nil, 0, err
	}
	return m, 1 + n, nil
}
