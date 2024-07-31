package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	TEXT         = lipgloss.Color("#313244")
	NEUTRAL      = lipgloss.Color("#bac2de")
	VALID        = lipgloss.Color("#a6e3a1")
	INVALID      = lipgloss.Color("#eb6f92")
	NOT_EDITABLE = lipgloss.Color("#ebbcba")
	BORDER       = lipgloss.Color("#89b4fa")
)

var base = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Foreground(TEXT).BorderForeground(BORDER)
var borderStyle = base.Foreground(BORDER)

type model struct {
	game    *Game
	playing bool
	loaded  bool
	spinner spinner.Model
	b0t     *bot
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}
	if m.playing {
		return updateGame(msg, m)
	}
	return m, nil
}

func updateGame(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.b0t = initBot(m.game.playingBoard, false)
			return m, m.b0t.move()
		default:
			m.game.handleMove(msg.String())
		}
	case spinner.TickMsg:
		if !m.game.loaded {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case nextBotMove:
		if m.b0t.ip <= len(m.b0t.solution) {
			mv := m.b0t.solution[m.b0t.ip]
			m.game.handleMove(mv)
			m.b0t.ip++
			return m, m.b0t.move()
		}

	case fillBoardCmd:
		for _, mv := range m.b0t.solution {
			if mv.value == 0 {
				continue
			}
			m.game.playingBoard[mv.r][mv.c].value = mv.value
		}
	}

	if m.game.checkWin() {
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	if m.playing {
		return gameView(m)
	}
	return ""
}
func gameView(m model) string {
	if !m.game.loaded {
		return lipgloss.JoinHorizontal(lipgloss.Center, "Loading ", m.spinner.View(), "\n\n")
	}
	board := m.game.playingBoard.View(m.game.r, m.game.c)
	width := lipgloss.Width(board)
	board = borderStyle.UnsetPadding().BorderStyle(lipgloss.NormalBorder()).Render(board)
	board = lipgloss.NewStyle().Width(width * 3).Align(lipgloss.Center).Render(board)
	return lipgloss.JoinVertical(lipgloss.Center, "Sudoku", board, "\n\n")
}

func (b Board) View(r, c int) string {
	s := ""
	for x := range 9 {
		for y := range 9 {
			if x%3 == 0 && x != 0 && y == 0 {
				s += borderStyle.Render("─────────┼───────────┼─────────")
				s += "\n"
			}
			if y%3 == 0 && y != 0 {
				s += borderStyle.Render("│")
			}

			if x == r && y == c {
				if b[x][y].value == 0 {
					s += base.Background(NEUTRAL).Render(" ")
				} else if b[x][y].editable {
					s += base.Background(VALID).Render(fmt.Sprint(b[x][y].value))
				} else {
					s += base.Background(INVALID).Render(fmt.Sprint(b[x][y].value))
				}
			} else {
				if b[x][y].value == 0 {
					s += base.Render(" ")
				} else if b[x][y].editable {
					s += base.Foreground(VALID).Render(fmt.Sprint(b[x][y].value))
				} else {
					s += base.Foreground(NOT_EDITABLE).Render(fmt.Sprint(b[x][y].value))
				}
			}
		}
		if x != 8 {
			s += "\n"
		}
	}
	return s
}
