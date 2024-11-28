package noxnet

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/opennox/libs/binenc"
)

func init() {
	RegisterMessage(&MsgInform{}, true)
}

type InformCode byte

const (
	InfoSpellErr       = InformCode(0)
	InfoSpellOK        = InformCode(1)
	InfoAbility        = InformCode(2)
	InfoPlayerTimeout  = InformCode(3)
	InfoFlagRetrieve   = InformCode(4)
	InfoFlagCapture    = InformCode(5)
	InfoFlagPickup     = InformCode(6)
	InfoFlagDrop       = InformCode(7)
	InfoFlagRespawn    = InformCode(8)
	InfoBallScore      = InformCode(9)
	InfoCrownPickup    = InformCode(10)
	InfoCrownDrop      = InformCode(11)
	InfoObserver       = InformCode(12)
	InfoPlayerState    = InformCode(13)
	InfoKill           = InformCode(14)
	InfoStringID       = InformCode(15)
	InfoWrongTeam      = InformCode(16)
	InfoClassChanged   = InformCode(17)
	InfoPlayerExit     = InformCode(18)
	InfoPlayerWarp     = InformCode(19)
	InfoPlayerSecret   = InformCode(20)
	InfoPlayerWarpWait = InformCode(21)
)

type MsgInform struct {
	Inform Inform
}

func (*MsgInform) NetOp() Op {
	return MSG_INFORM
}

func (m *MsgInform) EncodeSize() int {
	return 1 + m.Inform.EncodeSize()
}

func (m *MsgInform) Encode(data []byte) (int, error) {
	if len(data) < m.EncodeSize() {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(m.Inform.InformCode())
	n, err := m.Inform.Encode(data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

func (m *MsgInform) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	code := InformCode(data[0])
	switch code {
	case InfoSpellErr, InfoSpellOK, InfoAbility:
		m.Inform = &InformSpell{Code: code}
	case InfoPlayerTimeout, InfoFlagRetrieve,
		InfoPlayerExit, InfoPlayerWarp, InfoPlayerSecret:
		m.Inform = &InformPlayer{Code: code}
	case InfoFlagRespawn, InfoWrongTeam:
		m.Inform = &InformTeam{Code: code}
	case InfoObserver:
		m.Inform = &InformObserver{}
	case InfoPlayerState:
		m.Inform = &InformPlayerState{}
	case InfoFlagCapture, InfoFlagPickup, InfoFlagDrop,
		InfoCrownPickup, InfoCrownDrop, InfoBallScore:
		m.Inform = &InformTeamPlayer{Code: code}
	case InfoKill:
		m.Inform = &InformKill{}
	case InfoStringID:
		m.Inform = &InformStringID{}
	case InfoClassChanged:
		m.Inform = &InformClassChanged{}
	case InfoPlayerWarpWait:
		m.Inform = &InformWarpWait{}
	default:
		return 0, fmt.Errorf("unsupported inform code: %d", code)
	}
	n, err := m.Inform.Decode(data[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

type Inform interface {
	InformCode() InformCode
	Encoded
}

type InformPlayer struct {
	Code     InformCode
	PlayerID uint32
}

func (m *InformPlayer) InformCode() InformCode {
	return m.Code
}

func (*InformPlayer) EncodeSize() int {
	return 4
}

func (m *InformPlayer) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.PlayerID)
	return 4, nil
}

func (m *InformPlayer) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.PlayerID = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}

type InformTeam struct {
	Code   InformCode
	TeamID uint32
}

func (m *InformTeam) InformCode() InformCode {
	return m.Code
}

func (*InformTeam) EncodeSize() int {
	return 4
}

func (m *InformTeam) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.TeamID)
	return 4, nil
}

func (m *InformTeam) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.TeamID = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}

type InformObserver struct {
	Prompt uint32
}

func (*InformObserver) InformCode() InformCode {
	return InfoObserver
}

func (*InformObserver) EncodeSize() int {
	return 4
}

func (m *InformObserver) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.Prompt)
	return 4, nil
}

func (m *InformObserver) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Prompt = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}

type InformPlayerState struct {
	State uint32
}

func (*InformPlayerState) InformCode() InformCode {
	return InfoPlayerState
}

func (*InformPlayerState) EncodeSize() int {
	return 4
}

func (m *InformPlayerState) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.State)
	return 4, nil
}

func (m *InformPlayerState) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.State = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}

type InformSpell struct {
	Code InformCode
	Ind  uint32
}

func (m *InformSpell) InformCode() InformCode {
	return m.Code
}

func (*InformSpell) EncodeSize() int {
	return 4
}

func (m *InformSpell) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.Ind)
	return 4, nil
}

func (m *InformSpell) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Ind = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}

type InformTeamPlayer struct {
	Code     InformCode
	PlayerID uint32
	TeamID   uint32
}

func (m *InformTeamPlayer) InformCode() InformCode {
	return m.Code
}

func (*InformTeamPlayer) EncodeSize() int {
	return 8
}

func (m *InformTeamPlayer) Encode(data []byte) (int, error) {
	if len(data) < 8 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.PlayerID)
	binary.LittleEndian.PutUint32(data[4:8], m.TeamID)
	return 8, nil
}

func (m *InformTeamPlayer) Decode(data []byte) (int, error) {
	if len(data) < 8 {
		return 0, io.ErrUnexpectedEOF
	}
	m.PlayerID = binary.LittleEndian.Uint32(data[0:4])
	m.TeamID = binary.LittleEndian.Uint32(data[4:8])
	return 8, nil
}

type InformClassChanged struct {
}

func (*InformClassChanged) InformCode() InformCode {
	return InfoClassChanged
}

func (*InformClassChanged) EncodeSize() int {
	return 0
}

func (m *InformClassChanged) Encode(data []byte) (int, error) {
	return 0, nil
}

func (m *InformClassChanged) Decode(data []byte) (int, error) {
	return 0, nil
}

type InformStringID struct {
	IsSign   byte
	StringID string
}

func (*InformStringID) InformCode() InformCode {
	return InfoStringID
}

func (m *InformStringID) EncodeSize() int {
	return 1 + len(m.StringID) + 1
}

func (m *InformStringID) Encode(data []byte) (int, error) {
	if len(data) < m.EncodeSize() {
		return 0, io.ErrShortBuffer
	}
	data[0] = m.IsSign
	n := binenc.CStringSet0(data[1:], m.StringID)
	return 1 + n, nil
}

func (m *InformStringID) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	m.IsSign = data[0]
	m.StringID = binenc.CString(data[1:])
	return 1 + len(m.StringID) + 1, nil
}

type InformKill struct {
	Unk1 uint16
	Unk2 uint16
	Unk3 uint16
	Unk4 uint16
	Unk5 byte
}

func (*InformKill) InformCode() InformCode {
	return InfoKill
}

func (*InformKill) EncodeSize() int {
	return 9
}

func (m *InformKill) Encode(data []byte) (int, error) {
	if len(data) < 9 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:2], m.Unk1)
	binary.LittleEndian.PutUint16(data[2:4], m.Unk2)
	binary.LittleEndian.PutUint16(data[4:6], m.Unk3)
	binary.LittleEndian.PutUint16(data[6:8], m.Unk4)
	data[8] = m.Unk5
	return 9, nil
}

func (m *InformKill) Decode(data []byte) (int, error) {
	if len(data) < 9 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Unk1 = binary.LittleEndian.Uint16(data[0:2])
	m.Unk2 = binary.LittleEndian.Uint16(data[2:4])
	m.Unk3 = binary.LittleEndian.Uint16(data[4:6])
	m.Unk4 = binary.LittleEndian.Uint16(data[6:8])
	m.Unk5 = data[8]
	return 9, nil
}

type InformWarpWait struct {
	Stage uint32
}

func (*InformWarpWait) InformCode() InformCode {
	return InfoPlayerWarpWait
}

func (*InformWarpWait) EncodeSize() int {
	return 4
}

func (m *InformWarpWait) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], m.Stage)
	return 4, nil
}

func (m *InformWarpWait) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Stage = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}
