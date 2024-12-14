package noxnet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgSeqImportant{}, true)
}

var _ netmsg.ComplexMessage = (*MsgSeqImportant)(nil)

type MsgSeqImportant struct {
	ID  uint16
	Msg netmsg.Message
}

func (*MsgSeqImportant) NetOp() netmsg.Op {
	return netmsg.MSG_SEQ_IMPORTANT
}

func (m *MsgSeqImportant) EncodeSize() int {
	return m.EncodeSizeWith(nil)
}

func (m *MsgSeqImportant) EncodeSizeWith(s *netmsg.State) int {
	n := s.EncodeSize(m.Msg)
	if n > 0xff {
		n = 0xff
	}
	return 3 + n
}

func (m *MsgSeqImportant) Encode(data []byte) (int, error) {
	return m.EncodeWith(nil, data)
}

func (m *MsgSeqImportant) EncodeWith(s *netmsg.State, data []byte) (int, error) {
	sz := s.EncodeSize(m.Msg)
	if sz > 0xff {
		return 0, errors.New("message is too large")
	}
	if len(data) < 3+sz {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:2], m.ID)
	data[2] = byte(sz)
	n, err := s.Encode(data[3:], m.Msg)
	if err != nil {
		return 0, err
	}
	return 3 + n, nil
}

func (m *MsgSeqImportant) Decode(data []byte) (int, error) {
	return m.DecodeWith(nil, data)
}

func (m *MsgSeqImportant) DecodeWith(s *netmsg.State, data []byte) (int, error) {
	if len(data) < 3 {
		return 0, io.ErrUnexpectedEOF
	}
	m.ID = binary.LittleEndian.Uint16(data[0:2])
	sz := int(data[2])
	if len(data) < 3+sz {
		return 0, io.ErrUnexpectedEOF
	}
	data = data[3 : 3+sz]
	m.Msg = nil
	p, n, err := s.DecodeNext(data)
	if err != nil {
		return 3 + sz, err
	}
	m.Msg = p
	if len(data[n:]) > 0 {
		return 0, fmt.Errorf("partial message decoded: %v", p.NetOp())
	}
	return 3 + sz, nil
}
