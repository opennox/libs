package noxsave

import (
	"encoding/binary"
	"errors"
	"io"

	crypt "github.com/opennox/noxcrypt"
)

func ReadFileInfo(f io.Reader) (*FileInfo, error) {
	r, err := NewReader(f)
	if err != nil {
		return nil, err
	}
	for {
		raw, err := r.ReadRawSection()
		if err == io.EOF {
			return nil, errors.New("no save file info")
		} else if err != nil {
			return nil, err
		}
		if raw.ID == SectFileInfo {
			var info FileInfo
			_, err := info.Decode(raw.Data)
			if err != nil {
				return nil, err
			}
			return &info, nil
		}
	}
}

type Reader struct {
	cr  *crypt.Reader
	r   io.Reader
	reg *Registry
}

func NewReader(r io.Reader) (*Reader, error) {
	cr, err := crypt.NewReader(r, crypt.SaveKey)
	if err != nil {
		return nil, err
	}
	rd := &Reader{cr: cr, r: cr, reg: registry}
	return rd, nil
}

func (r *Reader) SetRegistry(reg *Registry) {
	r.reg = reg
}

func (r *Reader) ReadRawSection() (*RawSection, error) {
	var buf [4]byte
	if _, err := r.cr.Read(buf[:4]); err != nil {
		return nil, err
	}
	id := SectionID(binary.LittleEndian.Uint32(buf[:4]))
	if _, err := r.cr.ReadAligned(buf[:4]); err != nil {
		return nil, err
	}
	sz := binary.LittleEndian.Uint32(buf[:4])
	data := make([]byte, sz)
	if _, err := io.ReadFull(r.cr, data); err != nil {
		return nil, err
	}
	return &RawSection{ID: id, Data: data}, nil
}

func (r *Reader) ReadRawSections() ([]*RawSection, error) {
	var out []*RawSection
	for {
		s, err := r.ReadRawSection()
		if err == io.EOF {
			return out, nil
		} else if err != nil {
			return out, err
		}
		out = append(out, s)
	}
}

func (r *Reader) ReadSection() (Section, error) {
	raw, err := r.ReadRawSection()
	if err != nil {
		return nil, err
	}
	return raw.DecodeWith(r.reg)
}

func (r *Reader) ReadSections() ([]Section, error) {
	var out []Section
	for {
		s, err := r.ReadSection()
		if err == io.EOF {
			return out, nil
		} else if err != nil {
			return out, err
		}
		out = append(out, s)
	}
}
