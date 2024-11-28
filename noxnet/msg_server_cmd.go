package noxnet

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/opennox/libs/binenc"
)

func init() {
	RegisterMessage(&MsgServerCmd{}, true)
}

type MsgServerCmd struct {
	ID      byte
	NetCode uint16
	Cmd     string
}

func (*MsgServerCmd) NetOp() Op {
	return MSG_SERVER_CMD
}

func (p *MsgServerCmd) EncodeSize() int {
	return 4 + 2*(len(p.Cmd)+1)
}

func (p *MsgServerCmd) Encode(data []byte) (int, error) {
	if len(p.Cmd) > 0xff-1 {
		return 0, errors.New("command is too long")
	}
	if len(data) < p.EncodeSize() {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.ID
	binary.LittleEndian.PutUint16(data[1:3], p.NetCode)
	data[3] = byte(len(p.Cmd) + 1)
	n := binenc.CStringSet16(data[4:], p.Cmd)
	data[4+n+0] = 0
	data[4+n+1] = 0
	return 4 + n + 2, nil
}

func (p *MsgServerCmd) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	p.ID = data[0]
	p.NetCode = binary.LittleEndian.Uint16(data[1:3])
	sz := int(data[3])
	if len(data) < 4+2*sz {
		return 0, io.ErrUnexpectedEOF
	}
	p.Cmd = binenc.CString16(data[4 : 4+2*sz])
	return 4 + 2*sz, nil
}
