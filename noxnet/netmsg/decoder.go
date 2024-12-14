package netmsg

import (
	"fmt"
	"io"
	"reflect"
	"slices"
)

func DecodeNext(data []byte) (Message, int, error) {
	return (*State).DecodeNext(nil, data)
}

func Decode(data []byte, p Message) (int, error) {
	return (*State).Decode(nil, data, p)
}

func (s *State) Decode(data []byte, p Message) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	if got, exp := Op(data[0]), p.NetOp(); got != exp {
		return 0, fmt.Errorf("expected packet: %v, got: %v", exp, got)
	}
	cp, ok := p.(ComplexMessage)
	if !ok {
		n, err := p.Decode(data[1:])
		if err != nil {
			return 0, err
		}
		return 1 + n, err
	}
	if s == nil {
		return 0, fmt.Errorf("decoder required for %T", cp)
	}
	n, err := cp.DecodeWith(s, data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

func (s *State) DecodeNext(data []byte) (Message, int, error) {
	if len(data) == 0 {
		return nil, 0, io.EOF
	}
	op := Op(data[0])
	rt, ok := byOp[op]
	if !ok {
		if s.IsClient {
			rt, ok = byOpCli[op]
		} else {
			rt, ok = byOpSrv[op]
		}
	}
	if !ok {
		n := op.Len()
		if n < 0 || 1+n >= len(data) {
			// unknown size
			return &Unknown{
				Op: op, Data: slices.Clone(data[1:]),
			}, len(data), nil
		}
		return &Unknown{
			Op:   op,
			Data: slices.Clone(data[1 : 1+n]),
		}, 1 + n, nil
	}
	p := reflect.New(rt).Interface().(Message)
	cp, ok := p.(ComplexMessage)
	if !ok {
		n, err := p.Decode(data[1:])
		if err != nil {
			return nil, 0, err
		}
		return p, 1 + n, nil
	}
	if s == nil {
		// unknown size
		return &Unknown{
			Op: op, Data: slices.Clone(data[1:]),
		}, len(data), nil
	}
	n, err := cp.DecodeWith(s, data[1:])
	if err != nil {
		return nil, 0, err
	}
	return cp, 1 + n, nil
}
