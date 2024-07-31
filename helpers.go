package main

import "math/rand/v2"

func exists[K comparable](arr []K, value K) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func nextPos(r, c int) (int, int) {
	if c == 8 {
		return r + 1, 0
	}
	return r, c + 1
}

func nextRandPos() (int, int) {
	return rand.IntN(9), rand.IntN(9)
}
func twoDto1(r, c int) int {
	return (r * 9) + c
}

func oneDto2(pos int) (int, int) {
	return (pos / 9), (pos % 9)
}

func isSame(first, second *Board) bool {
	for r := range 9 {
		for c := range 9 {
			if first[r][c].value != second[r][c].value {
				return false
			}
		}
	}
	return true
}

func (game *Game) checkWin() bool {
	for r := range 9 {
		for c := range 9 {
			if game.playingBoard[r][c].value == 0 {
				return false
			}
		}
	}
	return game.playingBoard.Valid()
}
