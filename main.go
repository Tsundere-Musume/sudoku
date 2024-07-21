package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Cell struct {
	value    int
	editable bool
}

type Coords struct {
	r, c int
}
type Board [9][9]Cell

type Sudoku struct {
	completeBoard Board
	playingBoard  Board
	cursorX       int
	cursorY       int
	timer         stopwatch.Model
	started       bool
	bot           *Bot
}

type Bot struct {
	botInstructions []Instruction
	board           Board
	ip              int
	timeout         time.Duration
}

type BotInstruction struct{}

func (s Sudoku) moveBot() tea.Msg {
	time.Sleep(s.bot.timeout)
	return BotMsg{}
}

func exists[K comparable](arr []K, value K) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func (board Board) isValid() bool {
	rows := [9][]int{}
	col := [9][]int{}
	grid := [9][]int{}

	for r := range 9 {
		for c := range 9 {
			num := board[r][c].value
			if num == 0 {
				continue
			}

			g := (3 * (r / 3)) + c/3
			if exists(rows[r], num) || exists(col[c], num) || exists(grid[g], num) {
				return false
			}
			rows[r] = append(rows[r], num)
			col[c] = append(col[c], num)
			grid[g] = append(grid[g], num)
		}
	}
	return true
}

func (s *Sudoku) checkWin() bool {
	for r := range 9 {
		for c := range 9 {
			if s.playingBoard[r][c].editable && s.playingBoard[r][c].value == 0 {
				return false
			}
		}
	}
	return s.playingBoard.isValid()
}

func (s *Sudoku) solve() {
	var backtrack func(int, int) bool
	backtrack = func(r int, c int) bool {
		if r == 9 {
			return true
		}

		newR := r
		newC := c
		if c == 8 {
			newR += 1
			newC = 0
		} else {
			newC += 1
		}
		if s.playingBoard[r][c].value != 0 {
			return backtrack(newR, newC)
		}
		tried := []int{}
		for num := rand.IntN(9) + 1; len(tried) < 9; num = rand.IntN(9) + 1 {
			if exists(tried, num) {
				continue
			}
			s.completeBoard[r][c].value = num
			s.playingBoard[r][c].value = num
			tried = append(tried, num)
			if s.completeBoard.isValid() && backtrack(newR, newC) {
				return true
			}
			s.completeBoard[r][c].value = 0
			s.playingBoard[r][c].value = 0
		}
		return false
	}
	backtrack(0, 0)
}

type Instruction struct {
	Coords
	value int
}

func (s *Sudoku) getSolution() []Instruction {
	solution := []Instruction{}
	var backtrack func(int, int) bool
	backtrack = func(r int, c int) bool {
		if r == 9 {
			return true
		}

		newR := r
		newC := c
		if c == 8 {
			newR += 1
			newC = 0
		} else {
			newC += 1
		}
		if !s.playingBoard[r][c].editable {
			return backtrack(newR, newC)
		}
		tried := []int{}
		original := s.playingBoard[r][c].value
		for num := rand.IntN(9) + 1; len(tried) < 9; num = rand.IntN(9) + 1 {
			if exists(tried, num) {
				continue
			}
			s.playingBoard[r][c].value = num
			solution = append(solution, Instruction{Coords: Coords{r: r, c: c}, value: num})
			tried = append(tried, num)
			//TODO: make a version of isValid that ignores editable cells
			if s.playingBoard.isValid() && backtrack(newR, newC) {
				s.playingBoard[r][c].value = original
				return true
			}
			s.playingBoard[r][c].value = original
			solution = append(solution, Instruction{Coords: Coords{r: r, c: c}, value: original})
		}
		return false
	}
	backtrack(0, 0)
	return solution
}

// TODO: complete this version
// func (s *Sudoku) rmElements(n int) {
// 	var backtrack func(int, int, int) bool
// 	backtrack = func(r int, c int, n int) bool {
// 		if n == 0 {
// 			return true
// 		}
// 		for {
// 			num := s.playingBoard[r][c].value
// 			for num == 0 && s.playingBoard[r][c].editable {
// 				r, c = rand.IntN(9), rand.IntN(9)
// 				num = s.playingBoard[r][c].value
// 			}
// 			s.playingBoard[r][c].value = 0
// 			s.playingBoard[r][c].editable = true
// 			if s.hasUniqueSolution() && backtrack(r, c, n-1) {
// 				return true
// 			}
// 			s.playingBoard[r][c].value = num
// 			s.playingBoard[r][c].editable = false
// 			return backtrack(r, c, n)
// 		}
// 	}
//
// 	backtrack(0, 0, n)
// }

func (s *Sudoku) removeElements(n int) {
	for i := 0; i < n; {
		r, c := rand.IntN(9), rand.IntN(9)
		num := s.playingBoard[r][c].value
		if num == 0 {
			continue
		}
		s.playingBoard[r][c].value = 0
		s.playingBoard[r][c].editable = true
		if s.hasUniqueSolution() {
			i++
			continue
		}
		s.playingBoard[r][c].value = num
		s.playingBoard[r][c].editable = false
	}
}
func (s *Sudoku) hasUniqueSolution() bool {
	var backtrack func(int, int) bool
	backtrack = func(r int, c int) bool {
		if r == 9 {
			return true
		}
		newR := r
		newC := c
		if c == 8 {
			newR += 1
			newC = 0
		} else {
			newC += 1
		}
		if (!s.playingBoard[r][c].editable) || (s.playingBoard[r][c].value != 0) {
			return backtrack(newR, newC)
		}
		tried := []int{}
		tried = append(tried, s.completeBoard[r][c].value)
		for num := rand.IntN(9) + 1; len(tried) < 9; num = rand.IntN(9) + 1 {
			if exists(tried, num) {
				continue
			}

			s.playingBoard[r][c].value = num
			tried = append(tried, num)
			if s.playingBoard.isValid() && backtrack(newR, newC) {
				return true
			}
			s.playingBoard[r][c].value = 0
		}
		return false
	}
	return !backtrack(0, 0)
}
func (b *Board) show() {
	for r, row := range b {
		for c, cell := range row {
			if r%3 == 0 && r != 0 && c == 0 {
				fmt.Println("---------------------")
			}
			if c%3 == 0 && c != 0 {
				fmt.Print("| ")
			}
			if cell.value == 0 {
				fmt.Print("  ")
			} else {
				fmt.Printf("%v ", cell.value)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func initGame() Sudoku {
	game := Sudoku{
		timer: stopwatch.NewWithInterval(time.Millisecond),
	}
	game.solve()
	game.removeElements(5)
	game.completeBoard.show()
	// game.playingBoard.show()
	return game
}

func (m Sudoku) Init() tea.Cmd {
	return m.timer.Init()
}

type BotMsg struct{}

func (m Sudoku) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h", "left":
			if m.cursorX > 0 {
				m.cursorX--
			}
		case "l", "right":
			if m.cursorX < 8 {
				m.cursorX++
			}
		case "j", "down":
			if m.cursorY < 8 {
				m.cursorY++
			}
		case "k", "up":
			if m.cursorY > 0 {
				m.cursorY--
			}
		case "a":
			a := m.getSolution()
			m.bot = &Bot{
				botInstructions: a,
				board:           m.playingBoard,
				timeout:         50 * time.Millisecond,
			}
			return m, m.moveBot
			// return m, tickCmd()
		default:
			num, err := strconv.Atoi(msg.String())
			if err != nil || num > 9 || num < 1 {
				break
			}
			if m.playingBoard[m.cursorY][m.cursorX].editable {
				m.playingBoard[m.cursorY][m.cursorX].value = num
			}
		}

	case BotMsg:
		if m.bot.ip == len(m.bot.botInstructions) {
			return m, nil
		}
		ins := m.bot.botInstructions[m.bot.ip]
		m.cursorY = ins.r
		m.cursorX = ins.c
		m.playingBoard[ins.r][ins.c].value = ins.value
		m.bot.ip++
		return m, m.moveBot
	}

	if m.checkWin() {
		return m, tea.Quit
	}
	if !m.started {
		m.started = true
		return m, m.timer.Start()
	}
	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)
	return m, cmd
}

var base = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
var style = base.Foreground(lipgloss.Color("#f5e0dc"))
var borderStyle = style.Foreground(lipgloss.Color("#89b4fa"))
var border = lipgloss.NewStyle().BorderForeground(lipgloss.Color("#89b4fa"))

const (
	NEUTRAL = lipgloss.Color("#bac2de")
	VALID   = lipgloss.Color("#a6e3a1")
	INVALID = lipgloss.Color("#f38ba8")
)

func (m Sudoku) View() string {
	s := ""
	for r := range 9 {
		for c := range 9 {
			if r%3 == 0 && r != 0 && c == 0 {
				s += borderStyle.Render("─────────┼───────────┼─────────")
				s += "\n"
			}
			if c%3 == 0 && c != 0 {
				s += borderStyle.Render("│")
			}
			if r == m.cursorY && c == m.cursorX {
				if m.playingBoard[r][c].value == 0 {
					s += style.Background(NEUTRAL).Render(" ")
				} else if m.playingBoard[r][c].editable {
					s += style.Foreground(lipgloss.Color("#313244")).Background(VALID).Render(fmt.Sprintf("%v", m.playingBoard[r][c].value))
				} else {
					s += style.Foreground(lipgloss.Color("#313244")).Background(INVALID).Render(fmt.Sprintf("%v", m.playingBoard[r][c].value))
				}
			} else {
				if m.playingBoard[r][c].value == 0 {
					s += style.Render(" ")
				} else if m.playingBoard[r][c].editable {
					s += style.Foreground(VALID).Render(fmt.Sprintf("%v", m.playingBoard[r][c].value))
				} else {
					s += style.Render(fmt.Sprintf("%v", m.playingBoard[r][c].value))
				}
			}
		}
		if r != 8 {
			s += "\n"
		}
	}
	width := lipgloss.Width(s)
	s = border.BorderStyle(lipgloss.NormalBorder()).Render(s)
	return lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.JoinVertical(lipgloss.Center, "Sudoku", lipgloss.NewStyle().Width(width*3).Align(lipgloss.Center).Render(s), "\n\n"), m.timer.View())
}
func main() {
	p := tea.NewProgram(initGame())
	if _, err := p.Run(); err != nil {
		fmt.Println("There is an error")
		os.Exit(1)
	}
}
