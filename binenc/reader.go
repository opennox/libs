package binenc

import (
	"bytes"
	"encoding/binary"
	"image"
	"io"
	"math"
	"unicode/utf16"

	"github.com/opennox/libs/types"
)

func NewReader(data []byte) *Reader {
	return &Reader{data: data}
}

type Reader struct {
	data []byte
	off  int
}

func (r *Reader) Reset(data []byte) {
	r.data = data
	r.off = 0
}

func (r *Reader) Offset() int {
	return r.off
}

func (r *Reader) Remaining() int {
	return len(r.data) - r.off
}

func (r *Reader) Err() error {
	if r.off >= len(r.data) {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (r *Reader) ReadNext(n int) ([]byte, bool) {
	if r.off+n > len(r.data) {
		return nil, false
	}
	v := r.data[r.off : r.off+n : r.off+n]
	r.off += n
	return v, true
}

func (r *Reader) ReadU8() (byte, bool) {
	if r.off+1 > len(r.data) {
		return 0, false
	}
	v := r.data[r.off]
	r.off++
	return v, true
}

func (r *Reader) ReadI8() (int8, bool) {
	if r.off+1 > len(r.data) {
		return 0, false
	}
	v := r.data[r.off]
	r.off++
	return int8(v), true
}

func (r *Reader) ReadU16() (uint16, bool) {
	if r.off+2 > len(r.data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint16(r.data[r.off:])
	r.off += 2
	return v, true
}

func (r *Reader) ReadI16() (int16, bool) {
	if r.off+2 > len(r.data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint16(r.data[r.off:])
	r.off += 2
	return int16(v), true
}

func (r *Reader) ReadU24() ([3]byte, bool) {
	if r.off+3 > len(r.data) {
		return [3]byte{}, false
	}
	var v [3]byte
	copy(v[:], r.data[r.off:])
	r.off += 3
	return v, true
}

func (r *Reader) ReadU32() (uint32, bool) {
	if r.off+4 > len(r.data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(r.data[r.off:])
	r.off += 4
	return v, true
}

func (r *Reader) ReadI32() (int32, bool) {
	if r.off+4 > len(r.data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(r.data[r.off:])
	r.off += 4
	return int32(v), true
}

func (r *Reader) ReadU64() (uint64, bool) {
	if r.off+8 > len(r.data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint64(r.data[r.off:])
	r.off += 8
	return v, true
}

func (r *Reader) ReadI64() (int64, bool) {
	if r.off+8 > len(r.data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint64(r.data[r.off:])
	r.off += 8
	return int64(v), true
}

func (r *Reader) ReadF32() (float32, bool) {
	if r.off+4 > len(r.data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(r.data[r.off:])
	r.off += 4
	return math.Float32frombits(v), true
}

func (r *Reader) ReadPointI32() (image.Point, bool) {
	if r.off+8 > len(r.data) {
		return image.Point{}, false
	}
	var v image.Point
	v.X = int(int32(binary.LittleEndian.Uint32(r.data[r.off+0:])))
	v.Y = int(int32(binary.LittleEndian.Uint32(r.data[r.off+4:])))
	r.off += 8
	return v, true
}

func (r *Reader) ReadPointF32() (types.Pointf, bool) {
	if r.off+8 > len(r.data) {
		return types.Pointf{}, false
	}
	var v types.Pointf
	v.X = math.Float32frombits(binary.LittleEndian.Uint32(r.data[r.off+0:]))
	v.Y = math.Float32frombits(binary.LittleEndian.Uint32(r.data[r.off+4:]))
	r.off += 8
	return v, true
}

func (r *Reader) Align(off int) bool {
	_, ok := r.ReadNext(8 - (r.off+off)%8)
	return ok
}

func (r *Reader) ReadAllBytes() []byte {
	b, _ := r.ReadBytes(r.Remaining())
	return b
}

func (r *Reader) ReadBytes(sz int) ([]byte, bool) {
	if sz == 0 {
		return nil, true
	}
	s, ok := r.ReadNext(sz)
	if !ok {
		return nil, false
	}
	b := make([]byte, sz)
	copy(b, s)
	return b, true
}

func (r *Reader) ReadString(sz int) (string, bool) {
	if r.off+sz > len(r.data) {
		return "", false
	}
	s := r.data[r.off : r.off+sz]
	r.off += sz
	if i := bytes.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	return string(s), true
}

func (r *Reader) ReadBytes8() ([]byte, bool) {
	sz, ok := r.ReadU8()
	if !ok {
		return nil, false
	}
	return r.ReadBytes(int(sz))
}

func (r *Reader) ReadBytes16() ([]byte, bool) {
	sz, ok := r.ReadU16()
	if !ok {
		return nil, false
	}
	return r.ReadBytes(int(sz))
}

func (r *Reader) ReadBytes32() ([]byte, bool) {
	sz, ok := r.ReadU32()
	if !ok {
		return nil, false
	}
	return r.ReadBytes(int(sz))
}

func (r *Reader) ReadString8() (string, bool) {
	sz, ok := r.ReadU8()
	if !ok {
		return "", false
	}
	return r.ReadString(int(sz))
}

func (r *Reader) ReadString16() (string, bool) {
	sz, ok := r.ReadU16()
	if !ok {
		return "", false
	}
	return r.ReadString(int(sz))
}

func (r *Reader) ReadString32() (string, bool) {
	sz, ok := r.ReadU32()
	if !ok {
		return "", false
	}
	return r.ReadString(int(sz))
}

func (r *Reader) ReadWString(sz int) (string, bool) {
	if r.off+2*sz > len(r.data) {
		return "", false
	}
	s := r.data[r.off : r.off+2*sz]
	r.off += 2 * sz
	sbuf := make([]uint16, sz)
	for i := range sbuf {
		c := binary.LittleEndian.Uint16(s)
		if c == 0 {
			sbuf = sbuf[:i]
			break
		}
		sbuf[i] = c
		s = s[2:]
	}
	return string(utf16.Decode(sbuf)), true
}

func (r *Reader) ReadWString8() (string, bool) {
	sz, ok := r.ReadU8()
	if !ok {
		return "", false
	}
	return r.ReadWString(int(sz))
}

func (r *Reader) ReadWString16() (string, bool) {
	sz, ok := r.ReadU16()
	if !ok {
		return "", false
	}
	return r.ReadWString(int(sz))
}

func (r *Reader) ReadWString32() (string, bool) {
	sz, ok := r.ReadU32()
	if !ok {
		return "", false
	}
	return r.ReadString(int(sz))
}
