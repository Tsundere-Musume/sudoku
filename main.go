package main

// TODO: FIX BOARD GENERATION
// TODO: FIX BOT FREEZE IF EDITABLE CELLS HAVE INVALID VALUE

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
	noTimeout     bool
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

func twoDto1(r, c int) int {
	return (r * 9) + c
}
func (s *Sudoku) rmElements(n int) {
	var backtrack func(int, int, int) bool
	backtrack = func(r int, c int, count int) bool {
		if count == n {
			return true
		}
		tried := []int{}
		for newR, newC := rand.IntN(9), rand.IntN(9); len(tried) < 81; newR, newC = rand.IntN(9), rand.IntN(9) {
			num := s.playingBoard[r][c].value
			if num != 0 {
				s.playingBoard[r][c].value = 0
				s.playingBoard[r][c].editable = true
				if s.has() && backtrack(newR, newC, count+1) {
					return true
				}
				s.playingBoard[r][c].value = num
				s.playingBoard[r][c].editable = false
			}
			if n := twoDto1(r, c); !exists(tried, n) {
				tried = append(tried, twoDto1(r, c))
			}
			r = newR
			c = newC
		}
		return false
	}

	fmt.Println("solvabel", backtrack(rand.IntN(9), rand.IntN(9), 0))
}

func (s *Sudoku) removeElements(n int) {
	for i := 0; i < n; {
		r, c := rand.IntN(9), rand.IntN(9)
		num := s.playingBoard[r][c].value
		if num == 0 {
			continue
		}
		s.playingBoard[r][c].value = 0
		s.playingBoard[r][c].editable = true
		if s.has() {
			i++
			continue
		}
		s.playingBoard[r][c].value = num
		s.playingBoard[r][c].editable = false
	}
}
func (s *Sudoku) has() bool {
	var backtrack func(int, int) bool
	backtrack = func(r int, c int) bool {
		if r == 9 {
			if s.isSame() {
				return false
			}
			return true
		}
		newR := r
		newC := c
		if c == 8 {
			newR++
			newC = 0
		} else {
			newC++
		}

		if s.playingBoard[r][c].value != 0 {
			return backtrack(newR, newC)
		}
		tried := make([]int, 0, 9)
		for num := rand.IntN(9) + 1; len(tried) < 9; num = rand.IntN(9) + 1 {
			// r := r
			// c := c
			if exists(tried, num) {
				continue
			}
			tried = append(tried, num)
			s.playingBoard[r][c].value = num
			if s.playingBoard.isValid() && backtrack(newR, newC) {
				s.playingBoard[r][c].value = 0
				return true
			}
			s.playingBoard[r][c].value = 0
		}
		return false
	}
	a := backtrack(0, 0)
	return !a
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
		tried := []int{s.completeBoard[r][c].value}
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
		timer:     stopwatch.NewWithInterval(time.Millisecond),
		noTimeout: true,
	}
	game.solve()
	// game.removeElements(40)
	game.rmElements(53)
	fmt.Println("Unique", game.has())
	// time.Sleep(time.Second * 1)
	game.completeBoard.show()
	// game.playingBoard.show()
	return game
}

func (m Sudoku) Init() tea.Cmd {
	return m.timer.Init()
}

type BotMsg struct{}

func (m Sudoku) isSame() bool {
	for r := range 9 {
		for c := range 9 {
			if m.playingBoard[r][c].value != m.completeBoard[r][c].value {
				return false
			}
		}
	}
	return true
}
func (m Sudoku) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.bot != nil && msg.String() != "q" && msg.String() != "ctrl+c" {
			break
		}
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
			if m.noTimeout {
				for _, ins := range a {
					m.playingBoard[ins.Coords.r][ins.Coords.c].value = ins.value
				}
				break
			}
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
		fmt.Println()
		fmt.Println()
		fmt.Println(m.isSame())
		fmt.Println()
		fmt.Println()
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
