//go:build js && wasm

package main

import "github.com/2asm/tic_tac_toe/tictactoe"

func main() {
	tictactoe.NewGame().Start()
}
