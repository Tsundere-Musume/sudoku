package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: FIX BOT FREEZE IF EDITABLE CELLS HAVE INVALID VALUE
// TODO: ADD DEBUG

func main() {
	game := new(Game)
	go func() {
		var mu sync.Mutex
		newGame(game)
		mu.Lock() // dont know if this is required
		game.loaded = true
		mu.Unlock()
	}()

	s := spinner.New()
	s.Spinner = spinner.Dot
	m := model{
		game:    game,
		playing: true,
		spinner: s,
		loaded:  true,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("There is an error")
		os.Exit(1)
	}
}
