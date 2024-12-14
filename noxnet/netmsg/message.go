package netmsg

import (
	"io"
	"reflect"

	"github.com/opennox/libs/binenc"
)

var (
	byOp    = make(map[Op]reflect.Type)
	byOpCli = make(map[Op]reflect.Type)
	byOpSrv = make(map[Op]reflect.Type)
)

func Register(p Message, dynamic bool) {
	op := p.NetOp()
	if _, ok := byOp[op]; ok {
		panic("already registered")
	}
	if _, ok := byOpCli[op]; ok {
		panic("already registered")
	}
	if _, ok := byOpSrv[op]; ok {
		panic("already registered")
	}
	byOp[op] = reflect.TypeOf(p).Elem()
	_, compl := p.(ComplexMessage)
	if dynamic || compl {
		opLen[op] = -1
	} else {
		opLen[op] = p.EncodeSize()
	}
}

func RegisterClient(p Message, dynamic bool) {
	op := p.NetOp()
	if _, ok := byOp[op]; ok {
		panic("already registered")
	}
	if _, ok := byOpCli[op]; ok {
		panic("already registered")
	}
	byOpCli[op] = reflect.TypeOf(p).Elem()
	opLen[op] = -1
}

func RegisterServer(p Message, dynamic bool) {
	op := p.NetOp()
	if _, ok := byOp[op]; ok {
		panic("already registered")
	}
	if _, ok := byOpSrv[op]; ok {
		panic("already registered")
	}
	byOpSrv[op] = reflect.TypeOf(p).Elem()
	opLen[op] = -1
}

type Message interface {
	NetOp() Op
	binenc.Encoded
}

type ComplexMessage interface {
	Message
	EncodeSizeWith(s *State) int
	EncodeWith(s *State, data []byte) (int, error)
	DecodeWith(s *State, data []byte) (int, error)
}

var _ Message = (*Unknown)(nil)

type Unknown struct {
	Op   Op
	Data []byte
}

func (p *Unknown) NetOp() Op {
	return p.Op
}

func (p *Unknown) EncodeSize() int {
	return len(p.Data)
}

func (p *Unknown) Encode(data []byte) (int, error) {
	if len(data) < len(p.Data) {
		return 0, io.ErrShortBuffer
	}
	n := copy(data, p.Data)
	return n, nil
}

func (p *Unknown) Decode(data []byte) (int, error) {
	panic("decoding unknown packet without op")
}
