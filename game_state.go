package main

import "strconv"

func (game *Game) moveTo(r, c int) {
	if r < 0 || c < 0 || r >= 9 || c >= 9 {
		return
	}
	game.r = r
	game.c = c
}

func (game *Game) handleMove(key string) {
	switch key {
	case "h", "left":
		game.c--
	case "l", "right":
		game.c++
	case "j", "down":
		game.r++
	case "k", "up":
		game.r--
	case " ":
		game.moveTo(oneDto2((twoDto1(game.r, game.c) + 1) % 81))
	default:
		num, err := strconv.Atoi(key)
		if err != nil || num > 9 || num < 1 {
			break
		}
		if game.playingBoard[game.r][game.c].editable {
			game.playingBoard[game.r][game.c].value = num
		}
		return
	}
	game.r = (game.r + 9) % 9
	game.c = (game.c + 9) % 9
}
