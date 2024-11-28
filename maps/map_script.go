package maps

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/opennox/noxscript/ns/asm"

	"github.com/opennox/opennox-lib/binenc"
)

func init() {
	RegisterSection(&Script{})
	RegisterSection(&ScriptData{})
}

func ReadScript(r io.Reader) (*Script, error) {
	d, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	return d.ReadScript()
}

func (r *Reader) ReadScript() (*Script, error) {
	data, err := r.readScriptRaw()
	if err != nil {
		return nil, err
	} else if len(data) == 0 {
		return nil, nil
	}
	r.m.Script = new(Script)
	err = r.m.Script.UnmarshalBinary(data)
	return r.m.Script, err
}

func (r *Reader) readScriptRaw() ([]byte, error) {
	for {
		sect, err := r.nextSection()
		if err == io.EOF {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		if sect == "ScriptObject" {
			return io.ReadAll(r.r)
		}
	}
}

type Script struct {
	Data []byte
}

func (*Script) MapSection() string {
	return "ScriptObject"
}

func (sect *Script) MarshalBinary() ([]byte, error) {
	data := make([]byte, 2+4+len(sect.Data))
	binary.LittleEndian.PutUint16(data[0:], 1) // version
	binary.LittleEndian.PutUint32(data[2:], uint32(len(sect.Data)))
	copy(data[6:], sect.Data)
	return data, nil
}

func (sect *Script) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

func (sect *Script) Decode(r *binenc.Reader) error {
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers != 1 {
		return fmt.Errorf("unsupported version of script section: %d", vers)
	}
	sect.Data, ok = r.ReadBytes32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (sect *Script) ReadScript() (*asm.Script, error) {
	if len(sect.Data) == 0 {
		return nil, nil
	}
	return asm.ReadScript(bytes.NewReader(sect.Data))
}

type ScriptData struct {
	Data []byte
}

func (*ScriptData) MapSection() string {
	return "ScriptData"
}

func (sect *ScriptData) MarshalBinary() ([]byte, error) {
	data := make([]byte, 2+1+len(sect.Data))
	binary.LittleEndian.PutUint16(data[0:], 1) // version
	if len(sect.Data) != 0 {
		data[2] = 1
	}
	copy(data[3:], sect.Data)
	return data, nil
}

func (sect *ScriptData) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

func (sect *ScriptData) Decode(r *binenc.Reader) error {
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers != 1 {
		return fmt.Errorf("unsupported version of script data section: %d", vers)
	}
	has, ok := r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if has == 0 {
		sect.Data = nil
		return nil
	}
	sect.Data = r.ReadAllBytes()
	return nil
}

type ScriptHandler struct {
	Ind  int16
	Func string
	Val3 uint32
}

func (w *ScriptHandler) EncodingSize() int {
	if w.Ind > 1 {
		return 2
	}
	return 6 + len(w.Func) + 4
}

func (w *ScriptHandler) MarshalBinary() ([]byte, error) {
	data := make([]byte, w.EncodingSize())
	if w.Ind > 1 {
		binary.LittleEndian.PutUint16(data[0:], uint16(w.Ind))
		return data, nil
	}
	i := 0
	binary.LittleEndian.PutUint16(data[i:], uint16(w.Ind))
	i += 2
	binary.LittleEndian.PutUint32(data[i:], uint32(len(w.Func)))
	i += 4
	i += copy(data[i:], w.Func)
	binary.LittleEndian.PutUint32(data[i:], w.Val3)
	i += 4
	return data, nil
}

func (w *ScriptHandler) Decode(r *binenc.Reader) error {
	*w = ScriptHandler{}
	var ok bool
	w.Ind, ok = r.ReadI16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	if w.Ind <= 1 {
		w.Func, ok = r.ReadString32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
		w.Val3, ok = r.ReadU32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	}
	return nil
}

func (w *ScriptHandler) UnmarshalBinary(data []byte) error {
	return w.Decode(binenc.NewReader(data))
}
