package tictactoe

type player struct {
	username string
	symbol   symbol
	cells    []coord
}

func newPlayer(username string, symbol symbol) *player {
	return &player{username, symbol, []coord{}}
}

func initialPlayers(turn int) [2]*player {
	return [2]*player{
		newPlayer("Bot", symbol(turn)),
		newPlayer("You", symbol(turn^1)),
	}
}

func (p *player) containsCoords(cs ...coord) bool {
	count := 0
	for _, c := range cs {
		yes := 0
		for _, pc := range p.cells {
			if pc == c {
				yes = 1
			}
		}
		count += yes
	}
	return count == len(cs)
}

