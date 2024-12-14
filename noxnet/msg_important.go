package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.RegisterClient(&MsgImportantCli{}, false)
	netmsg.RegisterServer(&MsgImportantSrv{}, false)
	netmsg.RegisterClient(&MsgImportantAckCli{}, false)
	netmsg.RegisterServer(&MsgImportantAckSrv{}, false)
}

type MsgImportantCli struct {
}

func (*MsgImportantCli) NetOp() netmsg.Op {
	return netmsg.MSG_IMPORTANT
}

func (*MsgImportantCli) EncodeSize() int {
	return 0
}

func (m *MsgImportantCli) Encode(data []byte) (int, error) {
	return 0, nil
}

func (m *MsgImportantCli) Decode(data []byte) (int, error) {
	return 0, nil
}

type MsgImportantSrv struct {
	ID uint32
}

func (*MsgImportantSrv) NetOp() netmsg.Op {
	return netmsg.MSG_IMPORTANT
}

func (*MsgImportantSrv) EncodeSize() int {
	return 4
}

func (m *MsgImportantSrv) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.ID)
	return 4, nil
}

func (m *MsgImportantSrv) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.ID = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}

type MsgImportantAckCli struct {
	ID uint32
}

func (*MsgImportantAckCli) NetOp() netmsg.Op {
	return netmsg.MSG_IMPORTANT_ACK
}

func (*MsgImportantAckCli) EncodeSize() int {
	return 4
}

func (m *MsgImportantAckCli) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.ID)
	return 4, nil
}

func (m *MsgImportantAckCli) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.ID = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}

type MsgImportantAckSrv struct {
	TS Timestamp
}

func (*MsgImportantAckSrv) NetOp() netmsg.Op {
	return netmsg.MSG_IMPORTANT_ACK
}

func (*MsgImportantAckSrv) EncodeSize() int {
	return 4
}

func (m *MsgImportantAckSrv) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], uint32(m.TS))
	return 4, nil
}

func (m *MsgImportantAckSrv) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.TS = Timestamp(binary.LittleEndian.Uint32(data[0:4]))
	return 4, nil
}
