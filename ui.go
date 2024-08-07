package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tabID int

const (
	tabGame tabID = iota
	tabSettings
	tabCount
)

type model struct {
	game       *Game
	playing    bool
	loaded     bool
	spinner    spinner.Model
	currWindow tabID
	b0t        *bot
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
		if k == "tab" {
			m.currWindow = (m.currWindow + 1) % tabCount
		}

		if k == "r" && m.loaded && m.currWindow == tabGame {
			m.loaded = false
			game := &Game{
				loaded:  make(chan struct{}, 1),
				timer:   stopwatch.NewWithInterval(time.Millisecond),
				options: &settingsModel{settings: m.game.options.settings},
			}
			m.game = game
			go func() {
				newGame(game)
				game.loaded <- struct{}{}
			}()
			s := spinner.New()
			s.Spinner = spinner.Dot
			m.spinner = s
			cmds = append(cmds, m.spinner.Tick)
		}
	}
	mod, cmd := updateGame(msg, m)
	cmds = append(cmds, cmd)
	m = mod.(model)
	mod, cmd = updateSettings(msg, m)
	cmds = append(cmds, cmd)
	return mod, tea.Batch(cmds...)
}

func updateGame(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.b0t != nil || m.currWindow != tabGame {
			break
		}
		switch msg.String() {
		case "a":
			if !m.loaded {
				break
			}
			m.b0t = initBot(m.game.problemBoard, true)
			// return m, m.b0t.move()
			cmds = append(cmds, m.b0t.move())

		default:
			m.game.handleMove(msg.String())
		}
	case spinner.TickMsg:
		select {
		case <-m.game.loaded:
			m.b0t = nil
			m.loaded = true
			cmds = append(cmds, m.game.timer.Init(), m.game.timer.Start())
		default:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case stopwatch.TickMsg, stopwatch.StartStopMsg:
		m.game.timer, cmd = m.game.timer.Update(msg)
		return m, cmd

	case nextBotMove:
		if m.b0t.ip < len(m.b0t.solution) && m.loaded {
			mv := m.b0t.solution[m.b0t.ip]
			m.game.handleMove(mv)
			m.b0t.ip++
			// return m, m.b0t.move()
			cmds = append(cmds, m.b0t.move())
		}

	case fillBoardCmd:
		for _, mv := range m.b0t.solution {
			if !m.loaded {
				break
			}
			if mv.value == 0 {
				continue
			}
			m.game.playingBoard[mv.r][mv.c].value = mv.value
		}
	}

	if m.game.checkWin() {
		// fmt.Println(isSame(&m.game.playingBoard, &m.game.solvedBoard))
		return m, tea.Quit
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	tabs := m.renderTabs()
	var s string
	switch m.currWindow {
	case tabGame:
		s = gameView(m)
	case tabSettings:
		s = m.game.options.View()
	default:
		s = ""
	}
	window := lipgloss.NewStyle().Height(lipgloss.Height(tabs) - 2).Width(60).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(s)
	window = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(borderColor).BorderLeft(false).Render(window)
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs, window) + "\n\n"
}

func (m model) renderTabs() string {
	tabNames := []string{
		"Sudoku",
		"Settings",
	}
	style := lipgloss.NewStyle().Foreground(borderColor)
	tabStyles := []string{
		style.Render("╭") + "\n" + style.Render("│"),
	}
	for idx := range tabNames {
		if idx == int(m.currWindow) {
			tabStyles = append(tabStyles, activeTab.Render(tabNames[idx]))
			continue
		}
		tabStyles = append(tabStyles, tab.Render(tabNames[idx]))
	}
	row := lipgloss.JoinVertical(lipgloss.Right, tabStyles...)
	gap := tabGap.Render(strings.Repeat("\n", max(0, 30-lipgloss.Height(row)))) + "\n" + style.Render("╰")
	row = lipgloss.JoinVertical(lipgloss.Bottom, row, gap)
	return row
}
func gameView(m model) string {
	if !m.loaded {
		return lipgloss.JoinHorizontal(lipgloss.Center, "Loading ", m.spinner.View(), "\n\n")
	}
	board := m.game.playingBoard.View(m.game.r, m.game.c)
	board = lipgloss.JoinVertical(lipgloss.Center, "\n\nSudoku", board, "\n\n")
	dbugWindow := ""
	if m.game.options.settings.Debug {
		dbug := strings.Builder{}
		dbug.WriteString(fmt.Sprintf("Time elapsed: %v\n\n", m.game.timer.View()))
		dbug.WriteString(fmt.Sprintf("Board valid: %v\n\n", m.game.playingBoard.Valid()))
		dbug.WriteString(fmt.Sprintf("Board same: %v\n\n", isSame(&m.game.playingBoard, &m.game.solvedBoard)))
		dbugWindow = borderStyle.Width(lipgloss.Width(board)).Render(dbug.String())
	}
	return lipgloss.JoinVertical(lipgloss.Center, board, dbugWindow)
}

func (b Board) View(r, c int) string {
	s := ""
	for x := range 9 {
		for y := range 9 {
			if x%3 == 0 && x != 0 && y == 0 {
				s += borderStyle.Render("━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━━")
				s += "\n"
			}
			if y%3 == 0 && y != 0 {
				s += borderStyle.Render("┃")
			}

			if x == r && y == c {
				if b[x][y].value == 0 {
					s += base.Background(active).Render(" ")
				} else if b[x][y].editable {
					s += base.Background(valid).Render(fmt.Sprint(b[x][y].value))
				} else {
					s += base.Background(invalid).Render(fmt.Sprint(b[x][y].value))
				}
			} else {
				if b[x][y].value == 0 {
					s += base.Render(" ")
				} else if b[x][y].editable {
					s += base.Foreground(valid).Render(fmt.Sprint(b[x][y].value))
				} else {
					s += base.Foreground(NOT_EDITABLE).Render(fmt.Sprint(b[x][y].value))
				}
			}
		}
		if x != 8 {
			s += "\n"
		}
	}
	return boardBorder.Render(s)
}
