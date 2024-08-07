package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

// INFO: FIX BOT FREEZE IF EDITABLE CELLS HAVE INVALID VALUE:: FIXED THIS BUT DON'T LIKE THE METHOD USED
// TODO: ADD DEBUG

func main() {
	options, err := parseSettings("./settings.json")
	if err != nil {
		fmt.Println(err)
		options = &settings{}
	}
	game := &Game{
		loaded:  make(chan struct{}, 1),
		timer:   stopwatch.NewWithInterval(time.Millisecond),
		options: &settingsModel{settings: options},
	}
	go func() {
		newGame(game)
		game.loaded <- struct{}{}
	}()

	s := spinner.New()
	s.Spinner = spinner.Dot
	m := model{
		game:       game,
		playing:    true,
		spinner:    s,
		currWindow: tabGame,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("There is an error")
	}
	os.Exit(0)
}
