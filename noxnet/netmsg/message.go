package netmsg

import (
	"fmt"
	"io"
	"reflect"
	"slices"

	"github.com/opennox/libs/binenc"
)

var (
	byOp = make(map[Op]reflect.Type)
)

func Register(p Message, dynamic bool) {
	op := p.NetOp()
	if _, ok := byOp[op]; ok {
		panic("already registered")
	}
	byOp[op] = reflect.TypeOf(p).Elem()
	if dynamic {
		opLen[op] = -1
	} else {
		opLen[op] = p.EncodeSize()
	}
}

type Message interface {
	NetOp() Op
	binenc.Encoded
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

func EncodeSize(p Message) int {
	return 1 + p.EncodeSize()
}

func Encode(data []byte, p Message) (int, error) {
	if sz := EncodeSize(p); len(data) < sz {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(p.NetOp())
	n, err := p.Encode(data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

func Append(data []byte, p Message) ([]byte, error) {
	sz := EncodeSize(p)
	orig := data
	i := len(orig)
	data = append(data, make([]byte, sz)...)
	buf := data[i : i+sz]
	_, err := Encode(buf, p)
	if err != nil {
		return orig, err
	}
	return data, nil
}

func DecodeAny(data []byte) (Message, int, error) {
	if len(data) == 0 {
		return nil, 0, io.EOF
	}
	op := Op(data[0])
	rt, ok := byOp[op]
	if !ok {
		return &Unknown{
			Op: op, Data: slices.Clone(data[1:]),
		}, len(data), nil
	}
	p := reflect.New(rt).Interface().(Message)
	n, err := p.Decode(data[1:])
	return p, 1 + n, err
}

func Decode(data []byte, p Message) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	if got, exp := Op(data[0]), p.NetOp(); got != exp {
		return 0, fmt.Errorf("expected packet: %v, got: %v", exp, got)
	}
	n, err := p.Decode(data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}
