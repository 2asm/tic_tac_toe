package tictactoe

type symbol int

const (
	_X symbol = iota
	_O
)

func (s symbol) String() string {
	if s == _X {
		return "X"
	}
	return "O"
}
