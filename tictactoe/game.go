//go:build js && wasm

package tictactoe

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

type gameStatus int

const (
	_IN_PROGRESS gameStatus = iota
	_WON
	_LOST
	_DRAW
)

func (gs gameStatus) String() string {
	switch gs {
	case _IN_PROGRESS:
		return "Game in progress"
	case _WON:
		return "You won"
	case _LOST:
		return "You lost"
	case _DRAW:
		return "Draw"
	default:
		return "ERROR_STATE"
	}
}

/* game */
type game struct {
	players [2]*player
	turn    int // 0 or 1
}

func (g *game) gameStat() gameStatus {
	for _, cond := range winningConditions {
		if g.players[0].containsCoords(cond...) {
			return _LOST
		}
		if g.players[1].containsCoords(cond...) {
			return _WON
		}
	}
	if len(g.players[0].cells)+len(g.players[1].cells) == 9 {
		return _DRAW
	}
	return _IN_PROGRESS
}

func (g *game) ended() bool {
	return g.gameStat() != _IN_PROGRESS
}

func (g *game) reset() {
	turn := rand.Intn(2)
	g.players = initialPlayers(turn)
	g.turn = turn
}

func (g *game) getBot() *player {
	return g.players[0]
}

func (g *game) getHuman() *player {
	return g.players[1]
}

func NewGame() *game {
	turn := rand.Intn(2)
	return &game{
		players: initialPlayers(turn),
		turn:    turn,
	}
}

func (g *game) coordUsed(c coord) bool {
	return g.players[0].containsCoords(c) || g.players[1].containsCoords(c)
}

func (g *game) getEmptyCoord() coord {
	out := coord{0, 0}
	for {
		out.x = rand.Intn(3)
		out.y = rand.Intn(3)
		if !g.coordUsed(out) {
			break
		}
	}
	return out
}

func (g *game) Start() {
	g.renderInfo()
	for {
		select {
		case <-restartChan:
			g.reset()
			g.renderInfo()
			g.render()
		case c := <-playerMoveChan: // human turn
			g.render()
			if g.ended() || g.turn != 1 {
				break
			}
			g.getHuman().cells = append(g.getHuman().cells, c)
			g.turn ^= 1
			g.render()
		default: // bot turn
			g.render()
			time.Sleep(time.Millisecond * 200)
			if g.ended() || g.turn != 0 {
				break
			}
			move_coord := g.getEmptyCoord()
			g.getBot().cells = append(g.getBot().cells, move_coord)
			g.turn ^= 1
			g.render()
		}
		if g.ended() {
			g.renderResult()
		}
	}
}

/* render */
var idToCoord = map[string]coord{
	"00": {0, 0}, "01": {0, 1}, "02": {0, 2},
	"10": {1, 0}, "11": {1, 1}, "12": {1, 2},
	"20": {2, 0}, "21": {2, 1}, "22": {2, 2},
}

var coordToId = map[coord]string{
	{0, 0}: "00", {0, 1}: "01", {0, 2}: "02",
	{1, 0}: "10", {1, 1}: "11", {1, 2}: "12",
	{2, 0}: "20", {2, 1}: "21", {2, 2}: "22",
}

var (
	grid           = make([][]js.Value, 0)
	playerMoveChan = make(chan coord)
	restartChan    = make(chan struct{})
	restart        js.Value
	info           js.Value
	result         js.Value
)

func init() {
	for i := range 3 {
		gi := []js.Value{}
		for j := range 3 {
			e := js.Global().Get("document").Call("getElementById", coordToId[coord{i, j}])
			gi = append(gi, e)
		}
		grid = append(grid, gi)
	}

	result = js.Global().Get("document").Call("getElementById", "result")
	restart = js.Global().Get("document").Call("getElementById", "restart")
	info = js.Global().Get("document").Call("getElementById", "youbot")

	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		id := this.Get("id").String()
		playerMoveChan <- idToCoord[id]
		return nil
	})

	for i := range 3 {
		for j := range 3 {
			grid[i][j].Call("addEventListener", "click", cb)
		}
	}

	restart.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		restartChan <- struct{}{}
		return nil
	}))
}

func (g *game) renderInfo() {
	info.Set("innerText", fmt.Sprintf("bot(%v) vs you(%v)", g.getBot().symbol, g.getHuman().symbol))
}

func (g *game) renderResult() {
	result.Set("innerText", g.gameStat().String())
}

func (g *game) render() {
	result.Set("innerText", "") // reset result

	for _, c := range g.getHuman().cells { // fix human cells
		grid[c.x][c.y].Set("disabled", "true")
		grid[c.x][c.y].Set("style", "color:black;")
		grid[c.x][c.y].Set("innerText", g.getHuman().symbol.String())
	}
	for _, c := range g.getBot().cells { // fix bot cells
		grid[c.x][c.y].Set("disabled", "true")
		grid[c.x][c.y].Set("style", "color:black;")
		grid[c.x][c.y].Set("innerText", g.getBot().symbol.String())
	}
	for i := range 3 {
		for j := range 3 {
			if !g.coordUsed(coord{i, j}) {
				disabled := g.turn == 0 || g.ended()
				grid[i][j].Set("disabled", disabled)
				grid[i][j].Set("style", "color:transparent;")
			}
		}
	}
}
