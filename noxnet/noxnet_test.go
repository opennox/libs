package noxnet

import (
	"context"
	"log/slog"
	"net/netip"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/noxnet/discover"
	"github.com/opennox/libs/noxnet/udpconn"
)

const (
	resendTick = 20 * time.Millisecond
)

func newServerAndClient(t testing.TB, e Engine) (*Server, *Client) {
	log := slog.Default()
	srvC, cliC := udpconn.NewPipe(log, 10)
	srvC.Addr = netip.AddrPortFrom(srvC.Addr.Addr(), udpconn.DefaultPort)
	srvC.Debug = true
	cliC.Debug = true
	t.Cleanup(func() {
		_ = cliC.Close()
		_ = srvC.Close()
	})
	srv := NewServer(log, srvC, e, nil)
	t.Cleanup(srv.Close)
	cli := NewClient(log, cliC)
	t.Cleanup(cli.Close)
	return srv, cli
}

type testEngine struct {
	t testing.TB

	Info  discover.MsgServerInfo
	OnTry func(req *MsgServerTryJoin) error
	Pass  string
}

func (e *testEngine) ServerInfo(addr netip.AddrPort) *discover.MsgServerInfo {
	v := e.Info
	return &v
}

func (e *testEngine) PreJoin(addr netip.AddrPort, req *MsgServerTryJoin) error {
	if e.OnTry == nil {
		return ErrWrongVer
	}
	return e.OnTry(req)
}

func (e *testEngine) CheckPass(addr netip.AddrPort, pass string) error {
	if e.Pass == "" {
		return nil
	}
	if e.Pass != pass {
		return ErrWrongPassword
	}
	return nil
}

func (e *testEngine) Connect(addr netip.AddrPort) (Player, error) {
	//TODO implement me
	panic("implement me")
}

func TestDiscover(t *testing.T) {
	e := &testEngine{t: t, Info: discover.MsgServerInfo{
		PlayersCur: 3,
		PlayersMax: 250,
		MapName:    "testmap",
		ServerName: "TestServer",
	}}
	srv, cli := newServerAndClient(t, e)
	ctx, cancel := context.WithTimeout(context.Background(), resendTick)
	defer cancel()

	out := make(chan ServerInfoResp, 10)
	err := cli.Discover(ctx, 0, out)
	must.NoError(t, err)
	close(out)

	var got []ServerInfoResp
	for v := range out {
		v.Info.Token = 0
		got = append(got, v)
	}
	must.Eq(t, []ServerInfoResp{
		{Addr: srv.LocalAddr(), Info: e.Info},
	}, got)
}

func TestPreJoin(t *testing.T) {
	const pass = "1234"
	var got atomic.Pointer[MsgServerTryJoin]
	e := &testEngine{
		t: t,
		OnTry: func(req *MsgServerTryJoin) error {
			got.Store(req)
			return ErrPasswordRequired
		},
		Pass: pass,
	}
	srv, cli := newServerAndClient(t, e)
	ctx, cancel := context.WithTimeout(context.Background(), 5*resendTick)
	defer cancel()
	addr := srv.LocalAddr()

	exp := MsgServerTryJoin{}
	err := cli.TryJoin(ctx, addr, exp)
	must.ErrorIs(t, err, ErrPasswordRequired)
	must.Eq(t, &exp, got.Load())

	err = cli.TryPassword(ctx, pass)
	must.NoError(t, err)
}
