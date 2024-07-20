package main

import (
	"fmt"
	"math/rand/v2"
)

type Board [9][9]int

type Sudoku struct {
	completeBoard Board
	playingBoard  Board
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
			num := board[r][c]
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
		if s.playingBoard[r][c] != 0 {
			return backtrack(newR, newC)
		}
		tried := []int{}
		for num := rand.IntN(9) + 1; len(tried) < 9; num = rand.IntN(9) + 1 {
			if exists(tried, num) {
				continue
			}

			s.completeBoard[r][c] = num
			s.playingBoard[r][c] = num
			tried = append(tried, num)
			if s.completeBoard.isValid() && backtrack(newR, newC) {
				return true
			}
			s.completeBoard[r][c] = 0
			s.playingBoard[r][c] = 0
		}
		return false
	}
	backtrack(0, 0)
}

// TODO: complete this version
func (s *Sudoku) rmElements(n int) {
	var backtrack func(int, int, int) bool
	backtrack = func(r int, c int, n int) bool {
		if n == 0 {
			return true
		}
		num := s.playingBoard[r][c]
		var newR int
		var newC int
		for z := num; z == 0; {
			newR, newC = rand.IntN(9), rand.IntN(9)
			z = s.playingBoard[newR][newC]
		}
		if num != 0 {
			s.playingBoard[r][c] = 0
			if s.hasUniqueSolution() && backtrack(newR, newC, n-1) {
				return true
			}
			s.playingBoard[r][c] = num
		}
		return backtrack(newR, newC, n)
	}

	backtrack(0, 0, n)
}

func (s *Sudoku) removeElements(n int) {
	for i := 0; i < n; {
		r, c := rand.IntN(9), rand.IntN(9)
		num := s.playingBoard[r][c]
		if num == 0 {
			continue
		}
		s.playingBoard[r][c] = 0
		if s.hasUniqueSolution() {
			i++
			continue
		}
		s.playingBoard[r][c] = num
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
		if s.playingBoard[r][c] != 0 {
			return backtrack(newR, newC)
		}
		tried := []int{}
		tried = append(tried, s.completeBoard[r][c])
		for num := rand.IntN(9) + 1; len(tried) < 9; num = rand.IntN(9) + 1 {
			if exists(tried, num) {
				continue
			}

			s.playingBoard[r][c] = num
			tried = append(tried, num)
			if s.playingBoard.isValid() && backtrack(newR, newC) {
				return true
			}
			s.playingBoard[r][c] = 0
		}
		return false
	}
	return !backtrack(0, 0)
}
func (b *Board) show() {
	for r, row := range b {
		for c, value := range row {
			if r%3 == 0 && r != 0 && c == 0 {
				fmt.Println("---------------------")
			}
			if c%3 == 0 && c != 0 {
				fmt.Print("| ")
			}
			if value == 0 {
				fmt.Print("  ")
			} else {
				fmt.Printf("%v ", value)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func initGame() Sudoku {
	game := Sudoku{}
	game.solve()
	game.removeElements(53)
	return game
}

func main() {
	game := initGame()
	game.playingBoard.show()
	game.solve()
	game.playingBoard.show()

}
