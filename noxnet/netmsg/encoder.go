package netmsg

import "io"

func EncodeSize(p Message) int {
	return (*State).EncodeSize(nil, p)
}

func Encode(data []byte, p Message) (int, error) {
	return (*State).Encode(nil, data, p)
}

func Append(data []byte, p Message) ([]byte, error) {
	return (*State).Append(nil, data, p)
}

func (s *State) EncodeSize(p Message) int {
	cp, ok := p.(ComplexMessage)
	if !ok || s == nil {
		return 1 + p.EncodeSize()
	}
	return 1 + cp.EncodeSizeWith(s)
}

func (s *State) Encode(data []byte, p Message) (int, error) {
	if sz := s.EncodeSize(p); len(data) < sz {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(p.NetOp())
	cp, ok := p.(ComplexMessage)
	if !ok || s == nil {
		n, err := p.Encode(data[1:])
		if err != nil {
			return 0, err
		}
		return 1 + n, nil
	}
	n, err := cp.EncodeWith(s, data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

func (s *State) Append(data []byte, p Message) ([]byte, error) {
	sz := s.EncodeSize(p)
	orig := data
	i := len(orig)
	data = append(data, make([]byte, sz)...)
	buf := data[i : i+sz]
	_, err := s.Encode(buf, p)
	if err != nil {
		return orig, err
	}
	return data, nil
}
