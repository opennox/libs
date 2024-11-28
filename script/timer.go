package script

import (
	"time"

	ns4 "github.com/opennox/noxscript/ns/v4"
)

func Frames(num int) ns4.Duration {
	return ns4.Frames(num)
}

func Time(dt time.Duration) ns4.Duration {
	return ns4.Time(dt)
}

type Duration = ns4.Duration

type TimeSource = ns4.TimeSource

func NewTimers(src ns4.TimeSource) *Timers {
	return &Timers{
		src:    src,
		cur:    ns4.NowWithSource(src),
		active: make(map[uint]*Timer),
	}
}

type Timers struct {
	src    TimeSource
	last   uint
	cur    ns4.Duration
	active map[uint]*Timer
}

func (t *Timers) stopTimer(id uint) {
	delete(t.active, id)
}

func (t *Timers) SetTimer(d ns4.Duration, fnc func()) *Timer {
	t.last++
	tm := &Timer{t: t, id: t.last, at: d.Add(t.cur), fnc: fnc}
	t.active[tm.id] = tm
	return tm
}

func (t *Timers) Tick() {
	t.cur = ns4.NowWithSource(t.src)
	for _, tm := range t.active {
		if tm.at.LessOrEq(t.cur) {
			tm.fnc()
			tm.fnc = nil
			t.stopTimer(tm.id)
		}
	}
}

type Timer struct {
	t   *Timers
	id  uint
	at  ns4.Duration
	fnc func()
}

func (t *Timer) Stop() {
	if t.fnc != nil {
		t.fnc = nil
		t.t.stopTimer(t.id)
	}
}
