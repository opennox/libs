package noxnet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgPlayerInput{}, true)
	netmsg.Register(&MsgMouse{}, false)
}

type MsgPlayerInput struct {
	Inputs []PlayerInput
}

func (*MsgPlayerInput) NetOp() netmsg.Op {
	return netmsg.MSG_PLAYER_INPUT
}

func (m *MsgPlayerInput) EncodeSize() int {
	sz := 1
	for _, v := range m.Inputs {
		sz += 4 + v.EncodeSize()
	}
	return sz
}

func (m *MsgPlayerInput) Encode(data []byte) (int, error) {
	sz := m.EncodeSize()
	if len(data) < sz {
		return 0, io.ErrShortBuffer
	}
	psz := sz - 1
	if psz > 0xff {
		return 0, errors.New("too many inputs for one packet")
	}
	data[0] = byte(psz)
	i := 1
	for _, v := range m.Inputs {
		binary.LittleEndian.PutUint32(data[i:], uint32(v.CtrlCode()))
		i += 4
		n, err := v.Encode(data[i:])
		if err != nil {
			return 0, err
		}
		i += n
	}
	return sz, nil
}

func (m *MsgPlayerInput) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Inputs = nil
	psz := int(data[0])
	data = data[1 : 1+psz]
	for len(data) > 0 {
		code := CtrlCode(data[0])
		if len(data) < 4 {
			return 0, io.ErrUnexpectedEOF
		}
		data = data[4:]
		var inp PlayerInput
		switch code.DataSize() {
		case 0:
			inp = &PlayerInput0{Code: code}
		case 1:
			inp = &PlayerInput1{Code: code}
		case 4:
			inp = &PlayerInput4{Code: code}
		default:
			return 0, errors.New("unsupported input encoding")
		}
		n, err := inp.Decode(data)
		if err != nil {
			return 0, err
		}
		data = data[n:]
		m.Inputs = append(m.Inputs, inp)
	}
	return 1 + psz, nil
}

type MsgMouse struct {
	X, Y uint16
}

func (*MsgMouse) NetOp() netmsg.Op {
	return netmsg.MSG_MOUSE
}

func (*MsgMouse) EncodeSize() int {
	return 4
}

func (m *MsgMouse) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:2], m.X)
	binary.LittleEndian.PutUint16(data[2:4], m.Y)
	return 4, nil
}

func (m *MsgMouse) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	m.X = binary.LittleEndian.Uint16(data[0:2])
	m.Y = binary.LittleEndian.Uint16(data[2:4])
	return 4, nil
}

type PlayerInput interface {
	CtrlCode() CtrlCode
	binenc.Encoded
}

type PlayerInput0 struct {
	Code CtrlCode
}

func (p *PlayerInput0) CtrlCode() CtrlCode {
	return p.Code
}

func (*PlayerInput0) EncodeSize() int {
	return 0
}

func (*PlayerInput0) Encode(data []byte) (int, error) {
	return 0, nil
}

func (*PlayerInput0) Decode(data []byte) (int, error) {
	return 0, nil
}

type PlayerInput1 struct {
	Code CtrlCode
	Val  byte
}

func (p *PlayerInput1) CtrlCode() CtrlCode {
	return p.Code
}

func (*PlayerInput1) EncodeSize() int {
	return 1
}

func (p *PlayerInput1) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Val
	return 1, nil
}

func (p *PlayerInput1) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Val = data[0]
	return 1, nil
}

type PlayerInput4 struct {
	Code CtrlCode
	Val  uint32
}

func (p *PlayerInput4) CtrlCode() CtrlCode {
	return p.Code
}

func (*PlayerInput4) EncodeSize() int {
	return 4
}

func (p *PlayerInput4) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], p.Val)
	return 4, nil
}

func (p *PlayerInput4) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Val = binary.LittleEndian.Uint32(data[0:4])
	return 4, nil
}

const (
	CCOrientation            = CtrlCode(1)
	CCMoveForward            = CtrlCode(2)
	CCMoveBackward           = CtrlCode(3)
	CCMoveLeft               = CtrlCode(4)
	CCMoveRight              = CtrlCode(5)
	CCAction                 = CtrlCode(6)
	CCJump                   = CtrlCode(7)
	CCChat                   = CtrlCode(8)
	CCTeamChat               = CtrlCode(9)
	CCReadSpellbook          = CtrlCode(10)
	CCToggleConsole          = CtrlCode(11)
	CCIncreaseWindowSize     = CtrlCode(12)
	CCDecreaseWindowSize     = CtrlCode(13)
	CCIncreaseGamma          = CtrlCode(14)
	CCDecreaseGamma          = CtrlCode(15)
	CCQuit                   = CtrlCode(16)
	CCQuitMenu               = CtrlCode(17)
	CCReadMap                = CtrlCode(18)
	CCInventory              = CtrlCode(19)
	CCSpellGestureUp         = CtrlCode(20)
	CCSpellGestureDown       = CtrlCode(21)
	CCSpellGestureLeft       = CtrlCode(22)
	CCSpellGestureRight      = CtrlCode(23)
	CCSpellGestureUpperRight = CtrlCode(24)
	CCSpellGestureUpperLeft  = CtrlCode(25)
	CCSpellGestureLowerRight = CtrlCode(26)
	CCSpellGestureLowerLeft  = CtrlCode(27)
	CCSpellPatternEnd        = CtrlCode(28)
	CCCastQueuedSpell        = CtrlCode(29)
	CCCastMostRecentSpell    = CtrlCode(30)
	CCCastSpell1             = CtrlCode(31)
	CCCastSpell2             = CtrlCode(32)
	CCCastSpell3             = CtrlCode(33)
	CCCastSpell4             = CtrlCode(34)
	CCCastSpell5             = CtrlCode(35)
	CCMapZoomIn              = CtrlCode(36)
	CCMapZoomOut             = CtrlCode(37)
	CCNextWeapon             = CtrlCode(38)
	CCQuickHealthPotion      = CtrlCode(39)
	CCQuickManaPotion        = CtrlCode(40)
	CCQuickCurePoisonPotion  = CtrlCode(41)
	CCNextSpellSet           = CtrlCode(42)
	CCPreviousSpellSet       = CtrlCode(43)
	CCSelectSpellSet         = CtrlCode(44)
	CCBuildTrap              = CtrlCode(45)
	CCServerOptions          = CtrlCode(46)
	CCTaunt                  = CtrlCode(47)
	CCLaugh                  = CtrlCode(48)
	CCPoint                  = CtrlCode(49)
	CCInvertSpellTarget      = CtrlCode(50)
	CCToggleRank             = CtrlCode(51)
	CCToggleNetstat          = CtrlCode(52)
	CCToggleGUI              = CtrlCode(53)
	CCAutoSave               = CtrlCode(54)
	CCAutoLoad               = CtrlCode(55)
	CCScreenShot             = CtrlCode(56)
	ccMax                    = CtrlCode(57)
)

type CtrlCode byte

func (code CtrlCode) String() string {
	if name := ctrlCodes[code]; name != "" {
		return name
	}
	return fmt.Sprintf("CtrlCode(%d)", int(code))
}

func (code CtrlCode) DataSize() int {
	switch code {
	case CCOrientation:
		return 1
	case CCMoveForward, CCMoveBackward,
		CCMoveLeft, CCMoveRight:
		return 1
	case CCSpellPatternEnd:
		return 1
	case CCCastQueuedSpell:
		return 4
	case CCCastMostRecentSpell:
		return 4
	}
	return 0
}

var ctrlCodes = [ccMax]string{
	0:                        "Null",
	CCOrientation:            "Orientation",
	CCMoveForward:            "MoveForward",
	CCMoveBackward:           "MoveBackward",
	CCMoveLeft:               "MoveLeft",
	CCMoveRight:              "MoveRight",
	CCAction:                 "Action",
	CCJump:                   "Jump",
	CCChat:                   "Chat",
	CCTeamChat:               "TeamChat",
	CCReadSpellbook:          "ReadSpellbook",
	CCToggleConsole:          "ToggleConsole",
	CCIncreaseWindowSize:     "IncreaseWindowSize",
	CCDecreaseWindowSize:     "DecreaseWindowSize",
	CCIncreaseGamma:          "IncreaseGamma",
	CCDecreaseGamma:          "DecreaseGamma",
	CCQuit:                   "Quit",
	CCQuitMenu:               "QuitMenu",
	CCReadMap:                "ReadMap",
	CCInventory:              "Inventory",
	CCSpellGestureUp:         "SpellGestureUp",
	CCSpellGestureDown:       "SpellGestureDown",
	CCSpellGestureLeft:       "SpellGestureLeft",
	CCSpellGestureRight:      "SpellGestureRight",
	CCSpellGestureUpperRight: "SpellGestureUpperRight",
	CCSpellGestureUpperLeft:  "SpellGestureUpperLeft",
	CCSpellGestureLowerRight: "SpellGestureLowerRight",
	CCSpellGestureLowerLeft:  "SpellGestureLowerLeft",
	CCSpellPatternEnd:        "SpellPatternEnd",
	CCCastQueuedSpell:        "CastQueuedSpell",
	CCCastMostRecentSpell:    "CastMostRecentSpell",
	CCCastSpell1:             "CastSpell1",
	CCCastSpell2:             "CastSpell2",
	CCCastSpell3:             "CastSpell3",
	CCCastSpell4:             "CastSpell4",
	CCCastSpell5:             "CastSpell5",
	CCMapZoomIn:              "MapZoomIn",
	CCMapZoomOut:             "MapZoomOut",
	CCNextWeapon:             "NextWeapon",
	CCQuickHealthPotion:      "QuickHealthPotion",
	CCQuickManaPotion:        "QuickManaPotion",
	CCQuickCurePoisonPotion:  "QuickCurePoisonPotion",
	CCNextSpellSet:           "NextSpellSet",
	CCPreviousSpellSet:       "PreviousSpellSet",
	CCSelectSpellSet:         "SelectSpellSet",
	CCBuildTrap:              "BuildTrap",
	CCServerOptions:          "ServerOptions",
	CCTaunt:                  "Taunt",
	CCLaugh:                  "Laugh",
	CCPoint:                  "Point",
	CCInvertSpellTarget:      "InvertSpellTarget",
	CCToggleRank:             "ToggleRank",
	CCToggleNetstat:          "ToggleNetstat",
	CCToggleGUI:              "ToggleGUI",
	CCAutoSave:               "AutoSave",
	CCAutoLoad:               "AutoLoad",
	CCScreenShot:             "ScreenShot",
}
