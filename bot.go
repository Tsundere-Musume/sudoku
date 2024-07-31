package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type bot struct {
	ip       int
	delay    bool
	solution []move
	timeout  time.Duration
}

type nextBotMove struct{}

type fillBoardCmd struct{}

func initBot(board Board, delay bool) *bot {
	b := &bot{
		solution: board.getSolution(),
		delay:    delay,
		timeout:  time.Millisecond * 25,
	}
	return b
}

func (b *bot) move() tea.Cmd {
	if b.delay {
		return func() tea.Msg {
			time.Sleep(b.timeout)
			return nextBotMove{}
		}
	}
	return func() tea.Msg { return fillBoardCmd{} }
}
