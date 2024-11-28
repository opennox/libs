package maps

import (
	"errors"
	"fmt"
	"image"
	"io"

	crypt "github.com/opennox/noxcrypt"
)

type Writer struct {
	cw     *crypt.Writer
	crc    uint32
	crcOff int64
}

type WriterAt interface {
	io.Writer
	io.WriterAt
}

type Header struct {
	Magic uint32
	Offs  image.Point
}

func NewWriter(w WriterAt, h Header) (*Writer, error) {
	cw, err := crypt.NewWriter(w, crypt.MapKey)
	if err != nil {
		return nil, err
	}
	wr := &Writer{cw: cw}
	if err := wr.writeHeader(h); err != nil {
		return nil, err
	}
	return wr, nil
}

func (w *Writer) writeHeader(h Header) error {
	if h.Magic == 0 {
		h.Magic = Magic
	}
	if err := w.cw.WriteU32(h.Magic); err != nil {
		return err
	}
	switch h.Magic {
	case MagicOld:
		// nop
	case Magic:
		off, err := w.cw.WriteEmpty()
		if err != nil {
			return err
		}
		w.crcOff = off
	default:
		return fmt.Errorf("unsupported magic: 0x%x", h.Magic)
	}
	if err := w.cw.WriteU32(uint32(h.Offs.X)); err != nil {
		return err
	}
	if err := w.cw.WriteU32(uint32(h.Offs.Y)); err != nil {
		return err
	}
	return nil
}

func (w *Writer) Flush() error {
	return w.cw.Flush()
}

func (w *Writer) Close() error {
	if err := w.cw.Flush(); err != nil {
		return err
	}
	if w.crcOff != 0 {
		if err := w.cw.WriteU32At(w.cw.CRC(), w.crcOff); err != nil {
			return err
		}
	}
	return w.cw.Close()
}

func (w *Writer) writeSectionName(name string) error {
	if len(name)+1 >= 0xff {
		return errors.New("section name is too long")
	}
	buf := make([]byte, 1+len(name)+1)
	buf[0] = byte(len(name)) + 1
	i := copy(buf[1:], name) + 1
	buf[i] = 0
	if _, err := w.cw.Write(buf); err != nil {
		return err
	}
	return w.cw.Flush()
}

func (w *Writer) WriteRawSection(s RawSection) error {
	if err := w.writeSectionName(s.Name); err != nil {
		return err
	}
	if err := w.cw.WriteU64(uint64(len(s.Data))); err != nil {
		return err
	}
	if _, err := w.cw.Write(s.Data); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteSection(s Section) error {
	data, err := s.MarshalBinary()
	if err != nil {
		return err
	}
	return w.WriteRawSection(RawSection{
		Name: s.MapSection(),
		Data: data,
	})
}

func (w *Writer) WriteSections(arr []Section) error {
	for _, s := range arr {
		if err := w.WriteSection(s); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) WriteRawSections(arr []RawSection) error {
	for _, s := range arr {
		if err := w.WriteRawSection(s); err != nil {
			return err
		}
	}
	return nil
}
