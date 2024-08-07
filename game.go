package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/charmbracelet/bubbles/stopwatch"
)

type Game struct {
	playingBoard Board
	problemBoard Board //WARN: don't like how there are three boards for one game
	solvedBoard  Board
	timer        stopwatch.Model
	loaded       chan struct{}
	r, c         int
}

type Board [9][9]Cell

type Cell struct {
	value    int
	editable bool
}

func (board Board) Valid() bool {
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

func (board *Board) fill() {
	var backtrack func(int, int) bool
	backtrack = func(r int, c int) bool {
		if r == 9 {
			return true
		}

		if board[r][c].value != 0 {
			return backtrack(nextPos(r, c))
		}

		tried := []int{}
		for len(tried) < 9 {
			num := rand.IntN(9) + 1
			if !exists(tried, num) {
				tried = append(tried, num)
				board[r][c].value = num
				if board.Valid() && backtrack(nextPos(r, c)) {
					return true
				}
				board[r][c].value = 0
			}
		}
		return false
	}
	backtrack(0, 0)
}

type move struct {
	r, c, value int
}

func (board Board) getSolution() []move {
	solution := []move{}
	var backtrack func(int, int) bool
	backtrack = func(r int, c int) bool {
		if r == 9 {
			return true
		}

		if !board[r][c].editable {
			return backtrack(nextPos(r, c))
		}

		tried := []int{}
		original := board[r][c].value
		for len(tried) < 9 {
			num := rand.IntN(9) + 1
			if !exists(tried, num) {
				board[r][c].value = num
				solution = append(solution, move{r: r, c: c, value: num})
				tried = append(tried, num)
				if board.Valid() && backtrack(nextPos(r, c)) {
					return true
				}
				board[r][c].value = original
				solution = append(solution, move{r: r, c: c, value: original})
			}
		}
		return false
	}
	backtrack(0, 0)
	return solution
}

// board has to be fully filled with valid elements
func (board *Board) removeElements(n int) {
	//TODO: add checks, edit: dont remember which checks I was talking in this comment
	original := board.copy()
	var backtrack func(int, int, int) bool
	backtrack = func(r int, c int, count int) bool {
		if count == n {
			return true
		}
		tried := []int{}
		for len(tried) < 81 {
			r1, c1 := nextRandPos()
			num := board[r][c].value
			if num != 0 {
				board[r][c].value = 0
				board[r][c].editable = true
				if !board.hasOtherSolution(&original) && backtrack(r1, c1, count+1) {
					return true
				}
				board[r][c].value = num
				board[r][c].editable = false
			}
			if n := twoDto1(r, c); !exists(tried, n) {
				tried = append(tried, twoDto1(r, c))
			}
			r, c = r1, c1
		}
		return false
	}

	backtrack(rand.IntN(9), rand.IntN(9), 0)
}

func (board Board) hasOtherSolution(sol *Board) bool {
	var backtrack func(int, int) bool
	backtrack = func(r int, c int) bool {
		if r == 9 {
			return !isSame(&board, sol)
		}

		if !board[r][c].editable {
			return backtrack(nextPos(r, c))
		}

		tried := []int{}
		for len(tried) < 9 {
			num := rand.IntN(9) + 1
			if !exists(tried, num) {
				tried = append(tried, num)
				board[r][c].value = num
				if board.Valid() && backtrack(nextPos(r, c)) {
					return true
				}
				board[r][c].value = 0
			}
		}
		return false
	}
	return backtrack(0, 0)
}

func (board Board) copy() Board {
	return board
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

func newGame(g *Game) {
	g.solvedBoard.fill()
	g.playingBoard = g.solvedBoard.copy()
	g.playingBoard.removeElements(53)
	g.problemBoard = g.playingBoard.copy()
}
