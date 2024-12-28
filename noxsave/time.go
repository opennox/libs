package noxsave

import (
	"encoding/binary"
	"io"
	"time"
)

type SystemTime struct {
	Year         uint16
	Month        uint16
	DayOfWeek    uint16
	Day          uint16
	Hour         uint16
	Minute       uint16
	Second       uint16
	Milliseconds uint16
}

func (ts *SystemTime) EncodeSize() int {
	return 16
}

func (ts *SystemTime) Encode(data []byte) (int, error) {
	if len(data) < 16 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:2], ts.Year)
	binary.LittleEndian.PutUint16(data[2:4], ts.Month)
	binary.LittleEndian.PutUint16(data[4:6], ts.DayOfWeek)
	binary.LittleEndian.PutUint16(data[6:8], ts.Day)
	binary.LittleEndian.PutUint16(data[8:10], ts.Hour)
	binary.LittleEndian.PutUint16(data[10:12], ts.Minute)
	binary.LittleEndian.PutUint16(data[12:14], ts.Second)
	binary.LittleEndian.PutUint16(data[14:16], ts.Milliseconds)
	return 16, nil
}

func (ts *SystemTime) Decode(data []byte) (int, error) {
	if len(data) < 16 {
		return 0, io.ErrUnexpectedEOF
	}
	ts.Year = binary.LittleEndian.Uint16(data[0:2])
	ts.Month = binary.LittleEndian.Uint16(data[2:4])
	ts.DayOfWeek = binary.LittleEndian.Uint16(data[4:6])
	ts.Day = binary.LittleEndian.Uint16(data[6:8])
	ts.Hour = binary.LittleEndian.Uint16(data[8:10])
	ts.Minute = binary.LittleEndian.Uint16(data[10:12])
	ts.Second = binary.LittleEndian.Uint16(data[12:14])
	ts.Milliseconds = binary.LittleEndian.Uint16(data[14:16])
	return 16, nil
}

func (ts *SystemTime) Time() time.Time {
	if ts == nil {
		return time.Time{}
	}
	return time.Date(
		int(ts.Year), time.Month(ts.Month), int(ts.Day),
		int(ts.Hour), int(ts.Minute), int(ts.Second),
		int(ts.Milliseconds)*int(time.Millisecond),
		time.Local,
	)
}
