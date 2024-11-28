package mapv0_test

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/script"
	"github.com/opennox/libs/types"
)

func (g *testGame) Players() []script.Player {
	return g.players
}

func (g *testGame) HostPlayer() script.Player {
	return g.host
}

type playersPrint struct {
	g *testGame
}

func (p playersPrint) Print(text string) {
	for _, p := range p.g.players {
		p.Print(text)
	}
}

func (g *testGame) Global() script.Printer {
	return playersPrint{g: g}
}

func (g *testGame) BlindPlayers(v bool) {
	for _, p := range g.players {
		p.Blind(v)
	}
}

func (g *testGame) newPlayer(name string, unit script.Unit, host bool) *testPlayer {
	p := &testPlayer{
		g:    g,
		name: name,
		unit: unit,
	}
	g.players = append(g.players, p)
	if host {
		if g.host != nil {
			panic("already set")
		}
		g.host = p
	}
	return p
}

type testPlayer struct {
	script.BasePlayer
	g     *testGame
	name  string
	unit  script.Unit
	blind bool
	msg   string
}

func (p *testPlayer) Name() string {
	return p.name
}

func (p *testPlayer) String() string {
	return "Player(" + p.Name() + ")"
}

func (p *testPlayer) IsHost() bool {
	return p == p.g.host
}

func (p *testPlayer) GetObject() script.Object {
	return p.unit
}

func (p *testPlayer) Pos() types.Pointf {
	return p.unit.Pos()
}

func (p *testPlayer) SetPos(v types.Pointf) {
	p.unit.SetPos(v)
}

func (p *testPlayer) Unit() script.Unit {
	return p.unit
}

func (p *testPlayer) Print(text string) {
	p.msg = text
}

func (p *testPlayer) Blind(v bool) {
	p.blind = v
}

func TestPlayerName(t *testing.T) {
	g := newGame(t)
	g.newPlayer("Test", nil, true)

	g.Exec(`
	local p = Nox.Players.host

	if p.name ~= "Test" then
		error("invalid name field")
	end

	if p:__tostring() ~= "Player(Test)" then
		error("invalid string conversion")
	end
`)
}

func TestPlayersList(t *testing.T) {
	g := newGame(t)
	g.newPlayer("Test1", nil, false)
	g.newPlayer("Test2", nil, true)

	g.Exec(`
	local p1 = Nox.Players[1]
	local p2 = Nox.Players[2]
	local host = Nox.Players.host

	local players = Nox.Players()

	if p1.name ~= players[1].name then
		error("invalid player")
	end
	if p2.name ~= players[2].name then
		error("invalid player")
	end
	if p2.name ~= host.name then
		error("invalid player")
	end
	if not p2.host then
		error("invalid player")
	end
	if not host.host then
		error("invalid player")
	end
`)
}

func TestPlayersPrint(t *testing.T) {
	g := newGame(t)
	p1 := g.newPlayer("Test1", nil, false)
	p2 := g.newPlayer("Test2", nil, true)

	g.Exec(`
	local host = Nox.Players.host

	host:Print("foo")
`)
	must.EqOp(t, "", p1.msg)
	must.EqOp(t, "foo", p2.msg)

	g.Exec(`
	Nox.Players.Print("bar")
`)
	must.EqOp(t, "bar", p1.msg)
	must.EqOp(t, "bar", p2.msg)
}

func TestPlayersBlind(t *testing.T) {
	g := newGame(t)
	p1 := g.newPlayer("Test1", nil, false)
	p2 := g.newPlayer("Test2", nil, true)

	g.Exec(`
	host = Nox.Players.host

	host:Blind()
`)
	must.EqOp(t, false, p1.blind)
	must.EqOp(t, true, p2.blind)

	g.Exec(`
	host:Blind(false)
`)
	must.EqOp(t, false, p1.blind)
	must.EqOp(t, false, p2.blind)

	g.Exec(`
	host:Blind(true)
`)
	must.EqOp(t, false, p1.blind)
	must.EqOp(t, true, p2.blind)

	g.Exec(`
	Nox.Players.Blind()
`)
	must.EqOp(t, true, p1.blind)
	must.EqOp(t, true, p2.blind)

	g.Exec(`
	Nox.Players.Blind(false)
`)
	must.EqOp(t, false, p1.blind)
	must.EqOp(t, false, p2.blind)

	g.Exec(`
	Nox.Players.Blind(true)
`)
	must.EqOp(t, true, p1.blind)
	must.EqOp(t, true, p2.blind)
}

func TestPlayerUnit(t *testing.T) {
	g := newGame(t)
	v := g.newUnit("Player1", 1, 2)
	g.newPlayer("Test1", v, true)
	g.Exec(`
	local p = Nox.Players.host
	local v = p.unit
	if v:__tostring() ~= "Unit(Player1)" then
		error("invalid unit")
	end
`)
}

func TestPlayerUnitPos(t *testing.T) {
	g := newGame(t)
	v := g.newUnit("Player1", 1, 2)
	g.newPlayer("Test1", v, true)
	g.Exec(`
	local p = Nox.Players.host
	local v = p.unit

	if p.x ~= 1 then
		error("invalid X field")
	end
	if p.y ~= 2 then
		error("invalid Y field")
	end

	local x, y = p:Pos()
	if x ~= 1 then
		error("invalid X in Pos")
	end
	if y ~= 2 then
		error("invalid Y in Pos")
	end

	p.x, p.y = 3, 4
`)
	must.EqOp(t, types.Pointf{3, 4}, v.pos)
}

func TestPlayerPosArg(t *testing.T) {
	g := newGame(t)
	v1 := g.newUnit("Player1", 1, 2)
	v2 := g.newUnit("Unit2", 3, 4)
	g.newPlayer("Test1", v1, true)
	g.Exec(`
	local p = Nox.Players.host
	local v = Nox.Object("Unit2")
	v:SetPos(p)
`)
	must.EqOp(t, types.Pointf{1, 2}, v2.pos)
}

func TestPlayerObjArg(t *testing.T) {
	g := newGame(t)
	v1 := g.newUnit("Player1", 1, 2)
	v2 := g.newUnit("Unit2", 3, 4)
	g.newPlayer("Test1", v1, true)
	g.Exec(`
	local p = Nox.Players.host
	local v = Nox.Object("Unit2")
	v:Follow(p)
`)
	must.EqOp(t, UnitFollow, v2.st)
	must.EqOp(t, types.Pointf{1, 2}, v2.targ)
}
