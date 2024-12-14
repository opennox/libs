package noxnet

import (
	"io"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgAbilityAward{}, false)
}

type MsgAbilityAward struct {
	Ability byte
	Level   byte
}

func (*MsgAbilityAward) NetOp() netmsg.Op {
	return netmsg.MSG_REPORT_ABILITY_AWARD
}

func (*MsgAbilityAward) EncodeSize() int {
	return 2
}

func (p *MsgAbilityAward) Encode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Ability
	data[1] = p.Level
	return 2, nil
}

func (p *MsgAbilityAward) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Ability = data[0]
	p.Level = data[1]
	return 2, nil
}
