package noxnet

import (
	"errors"
	"fmt"
	"io"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgServerError{}, false)
	netmsg.Register(&MsgPasswordRequired{}, false)
	netmsg.Register(&MsgJoinFailed{}, false)
}

type ErrorMsg interface {
	netmsg.Message
	Error() error
}

func ErrorToMsg(err error) (ErrorMsg, bool) {
	if err == nil {
		return nil, true
	}
	if e := (*ConnectError)(nil); errors.As(err, &e) {
		return &MsgServerError{Err: *e}, true
	} else if errors.Is(err, ErrPasswordRequired) {
		return &MsgPasswordRequired{}, true
	}
	return nil, false
}

var _ error = ConnectError(0)

type ConnectError byte

func (e ConnectError) Name() string {
	switch e {
	case ErrLowPing:
		return "ErrLowPing"
	case ErrHighPing:
		return "ErrHighPing"
	case ErrLowLevel:
		return "ErrLowLevel"
	case ErrHighLevel:
		return "ErrHighLevel"
	case ErrClosed:
		return "ErrClosed"
	case ErrBanned:
		return "ErrBanned"
	case ErrWrongPassword:
		return "ErrWrongPassword"
	case ErrIllegalClass:
		return "ErrIllegalClass"
	case ErrTimeOut:
		return "ErrTimeOut"
	case ErrFindFailed:
		return "ErrFindFailed"
	case ErrNeedRefresh:
		return "ErrNeedRefresh"
	case ErrFull:
		return "ErrFull"
	case ErrDupSerial:
		return "ErrDupSerial"
	case ErrWrongVer:
		return "ErrWrongVer"
	}
	return fmt.Sprintf("ConnectError(%d)", int(e))
}

func (e ConnectError) Error() string {
	switch e {
	case ErrLowPing:
		return "ping is too low"
	case ErrHighPing:
		return "ping is too high"
	case ErrLowLevel:
		return "level is too low"
	case ErrHighLevel:
		return "level is too high"
	case ErrClosed:
		return "server is closed"
	case ErrBanned:
		return "banned on the server"
	case ErrWrongPassword:
		return "wrong password"
	case ErrIllegalClass:
		return "illegal player class"
	case ErrTimeOut:
		return "server timeout"
	case ErrFindFailed:
		return "find failed"
	case ErrNeedRefresh:
		return "needs refresh"
	case ErrFull:
		return "server is full"
	case ErrDupSerial:
		return "duplicate serial"
	case ErrWrongVer:
		return "wrong version"
	}
	return e.Name()
}

const (
	ErrLowPing = ConnectError(iota)
	ErrHighPing
	ErrLowLevel
	ErrHighLevel
	ErrClosed
	ErrBanned
	ErrWrongPassword
	ErrIllegalClass
	ErrTimeOut
	ErrFindFailed
	ErrNeedRefresh
	ErrFull
	ErrDupSerial
	ErrWrongVer
)

var _ ErrorMsg = (*MsgServerError)(nil)

type MsgServerError struct {
	Err ConnectError
}

func (p *MsgServerError) Error() error {
	return p.Err
}

func (*MsgServerError) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_ERROR
}

func (*MsgServerError) EncodeSize() int {
	return 1
}

func (p *MsgServerError) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = byte(p.Err)
	return 1, nil
}

func (p *MsgServerError) Decode(data []byte) (int, error) {
	if len(data) < 9 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Err = ConnectError(data[0])
	return 1, nil
}

var _ ErrorMsg = (*MsgPasswordRequired)(nil)

type MsgPasswordRequired struct {
}

func (*MsgPasswordRequired) Error() error {
	return ErrPasswordRequired
}

func (*MsgPasswordRequired) NetOp() netmsg.Op {
	return netmsg.MSG_PASSWORD_REQUIRED
}

func (*MsgPasswordRequired) EncodeSize() int {
	return 0
}

func (p *MsgPasswordRequired) Encode(data []byte) (int, error) {
	return 0, nil
}

func (p *MsgPasswordRequired) Decode(data []byte) (int, error) {
	return 0, nil
}

var _ ErrorMsg = (*MsgJoinFailed)(nil)

type MsgJoinFailed struct {
}

func (*MsgJoinFailed) Error() error {
	return ErrJoinFailed
}

func (*MsgJoinFailed) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_JOIN_FAIL
}

func (*MsgJoinFailed) EncodeSize() int {
	return 0
}

func (p *MsgJoinFailed) Encode(data []byte) (int, error) {
	return 0, nil
}

func (p *MsgJoinFailed) Decode(data []byte) (int, error) {
	return 0, nil
}
