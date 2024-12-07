package netxfer

type ReceiveFunc[C Conn] func(conn C, p Data)

type recvChunk struct {
	ind  Chunk
	data []byte
	next *recvChunk
	prev *recvChunk
}

type recvStream[C Conn] struct {
	conn       C
	recvID     RecvID
	full       Data
	nextChunk  int
	received   int
	progress   float64
	lastUpdate Timestamp
	first      *recvChunk
	last       *recvChunk
}

func (p *recvStream[C]) Reset(id RecvID) {
	*p = recvStream[C]{
		conn:      p.conn,
		recvID:    id,
		nextChunk: 1,
	}
}

func (p *recvStream[C]) AddChunk(chunk Chunk, data []byte) {
	if int(chunk) == p.nextChunk {
		copy(p.full.Data[p.received:p.received+len(data)], data)
		p.received += len(data)
		p.nextChunk++
	} else {
		b := &recvChunk{
			ind:  chunk,
			data: make([]byte, len(data)),
		}
		copy(b.data, data)
		b.next = nil
		b.prev = p.last
		if last := p.last; last != nil {
			last.next = b
		}
		p.last = b
		if first := p.first; first == nil {
			p.first = b
		}
	}

	for it := p.first; it != nil; it = it.next {
		if p.nextChunk == int(it.ind) {
			copy(p.full.Data[p.received:p.received+len(it.data)], it.data)
			p.received += len(it.data)
			p.nextChunk++
			if prev := it.prev; prev != nil {
				prev.next = it.next
			} else {
				p.first = it.next
			}
			if next := it.next; next != nil {
				next.prev = it.prev
			} else {
				p.last = it.prev
			}
			*it = recvChunk{}
		}
	}
	p.progress = float64(p.received) / float64(len(p.full.Data)) * 100.0
}

type Receiver[C Conn] struct {
	arr    []recvStream[C]
	cnt    int
	active int

	onReceive ReceiveFunc[C]
}

func (x *Receiver[C]) Init(n int, onRecv ReceiveFunc[C]) {
	if n < 0 {
		n = minStreams
	} else if n > maxStreams {
		n = maxStreams
	}
	x.onReceive = onRecv
	x.cnt = n
	x.active = 0
	x.arr = make([]recvStream[C], n)
	for i := 0; i < n; i++ {
		x.arr[i].Reset(RecvID(i))
	}
}

func (x *Receiver[C]) reset(rid RecvID) {
	x.arr[rid].Reset(rid)
}

func (x *Receiver[C]) Free() {
	for i := 0; i < x.cnt; i++ {
		x.free(RecvID(i))
	}
	x.arr = nil
}

func (x *Receiver[C]) free(rid RecvID) {
	p := &x.arr[rid]
	for it := p.first; it != nil; it = it.next {
		*it = recvChunk{}
	}
}

func (x *Receiver[C]) Update(ts Timestamp) {
	for i := 0; i < x.cnt; i++ {
		it := &x.arr[i]
		if len(it.full.Data) != 0 && ts > it.lastUpdate+900 {
			x.abort(RecvID(i), ErrRecvTimeout)
		}
	}
}

func (x *Receiver[C]) abort(rid RecvID, reason Error) {
	p := &x.arr[rid]
	if len(p.full.Data) != 0 {
		_ = p.conn.SendReliableMsg(&MsgAbort{
			RecvID: p.recvID,
			Reason: reason,
		})
		if x.active != 0 {
			x.active--
		}
		x.free(rid)
		x.reset(rid)
	}
}

func (x *Receiver[C]) newStream(conn C, act Action, typ string, sz uint32) *recvStream[C] {
	if x.cnt <= 0 {
		return nil
	}
	for id := 0; id < x.cnt; id++ {
		it := &x.arr[id]
		if len(it.full.Data) == 0 {
			it.recvID = RecvID(id)
			it.conn = conn
			it.full = Data{
				Action: act,
				Type:   typ,
				Data:   make([]byte, sz),
			}
			return it
		}
	}
	return nil
}

func (x *Receiver[C]) HandleStart(conn C, ts Timestamp, m *MsgStart) {
	if m.Size == 0 {
		return
	}
	p := x.newStream(conn, m.Act, m.Type.Value, m.Size)
	if p == nil {
		return
	}
	p.lastUpdate = ts
	x.active++
	_ = conn.SendReliableMsg(&MsgAccept{
		RecvID: p.recvID,
		SendID: m.SendID,
	})
}

func (x *Receiver[C]) HandleData(conn C, ts Timestamp, m *MsgData) {
	id := m.RecvID
	chunk := m.Chunk
	_ = conn.SendReliableMsg(&MsgAck{
		RecvID: id,
		Token:  m.Token,
		Chunk:  chunk,
	})
	if len(m.Data) == 0 {
		return
	}
	if int(id) >= x.cnt {
		return
	}
	p := &x.arr[id]
	if len(p.full.Data) == 0 {
		return
	}
	p.lastUpdate = ts
	p.AddChunk(chunk, m.Data)
	if p.received == len(p.full.Data) {
		_ = p.conn.SendReliableMsg(&MsgDone{
			RecvID: p.recvID,
		})
		if x.onReceive != nil {
			x.onReceive(conn, p.full)
		}
		if x.active != 0 {
			x.active--
		}
		p.full = Data{}
		x.free(id)
		x.reset(id)
	}
}

func (x *Receiver[C]) HandleCancel(m *MsgCancel) {
	id := m.RecvID
	x.free(id)
	x.reset(id)
}

func (x *Receiver[C]) Handle(conn C, ts Timestamp, m Msg) {
	switch m := m.(type) {
	case *MsgStart:
		x.HandleStart(conn, ts, m)
	case *MsgCancel:
		x.HandleCancel(m)
	case *MsgData:
		x.HandleData(conn, ts, m)
	}
}
