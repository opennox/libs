package udpconn

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/noxnet/discover"
	"github.com/opennox/libs/noxnet/netmsg"
)

func TestSeqBefore(t *testing.T) {
	cases := []struct {
		name   string
		v, cur Seq
		exp    bool
	}{
		{"zero", 0, 0, true},
		{"equal", 100, 100, true},
		{"future", 10, 5, false},
		{"max left", maxAckMsgs / 4, maxAckMsgs / 2, true},
		{"max right", 0xff - maxAckMsgs/2, 0xff - maxAckMsgs/4, true},
		{"overflow", 0xff - maxAckMsgs/4, maxAckMsgs / 4, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			must.EqOp(t, c.exp, c.v.Before(c.cur))
		})
	}
}

func TestStream(t *testing.T) {
	const maxMessages = 300

	newTest := func(t testing.TB, fast, debug bool, drop func(b []byte) bool) (*Conn, <-chan netmsg.Message) {
		srvC, cliC := NewPipe(slog.Default(), maxAckMsgs)
		t.Cleanup(func() {
			_ = cliC.Close()
			_ = srvC.Close()
		})
		cliC.Drop = drop
		cliC.Debug = debug
		srvC.Debug = debug
		srvAddr := srvC.Addr

		srvRecv := make(chan netmsg.Message, maxMessages)
		srv := srvC.Port
		srv.OnMessage(func(s Stream, m netmsg.Message, flags PacketFlags) bool {
			if debug {
				t.Logf("server recv: %#v", m)
			}
			if fast {
				_ = s.Conn().Ack()
			}
			srvRecv <- m
			return true
		})
		t.Cleanup(srv.Close)

		cli := cliC.Port
		cli.OnMessage(func(s Stream, m netmsg.Message, flags PacketFlags) bool {
			if debug {
				t.Logf("client recv: %#v", m)
			}
			return true
		})
		t.Cleanup(cli.Close)

		srv.Start()
		cli.Start()

		return cli.Conn(srvAddr), srvRecv
	}

	// Test that general ACK mechanism works for a long sequence of messages.
	t.Run("sequential", func(t *testing.T) {
		cliConn, srvRecv := newTest(t, true, true, nil)

		ctx, cancel := context.WithTimeout(context.Background(), 5*resendTick)
		defer cancel()

		timer := time.NewTimer(resendTick)

		for i := 0; i < maxMessages; i++ {
			exp := &discover.MsgDiscover{Token: uint32(i + 1)}
			err := cliConn.SendReliable(ctx, 0, exp)
			if err != nil {
				t.Fatal(err)
			}
			timer.Reset(resendTick)
			select {
			case m := <-srvRecv:
				must.Eq[netmsg.Message](t, exp, m)
			case <-timer.C:
				t.Fatal("expected a message")
			}
		}
	})

	// Test that large queue works.
	t.Run("long queue", func(t *testing.T) {
		cliConn, srvRecv := newTest(t, false, false, nil)

		ctx := context.Background()

		var expected []netmsg.Message
		for i := 0; i < maxMessages; i++ {
			exp := &discover.MsgDiscover{Token: uint32(i + 1)}
			cliConn.QueueReliable(0, Options{Context: ctx}, exp)
			expected = append(expected, exp)
		}
		for _, exp := range expected {
			select {
			case m := <-srvRecv:
				must.Eq[netmsg.Message](t, exp, m)
			case <-ctx.Done():
				t.Fatal("expected a message")
			}
		}
	})

	// Test redeliveries.
	t.Run("redelivery", func(t *testing.T) {
		dropped := 0
		cliConn, srvRecv := newTest(t, false, true, func(data []byte) bool {
			if dropped < resendRetries-1 {
				dropped++
				return true
			}
			return false
		})

		ctx := context.Background()

		exp := &discover.MsgDiscover{Token: 0x123}
		err := cliConn.SendReliable(ctx, 0, exp)
		if err != nil {
			t.Fatal(err)
		}
		select {
		case m := <-srvRecv:
			must.Eq[netmsg.Message](t, exp, m)
		default:
			t.Fatal("expected a message")
		}
	})
}
