// not valid tests just some functions to check if board the ui elements print correctly

package main

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func TestBoard(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Run("TestBoard", func(t *testing.T) {
		board := Board{}
		board.fill()

		tried := []int{}
		count := 0
		for r, c := nextRandPos(); count < 45; r, c = nextRandPos() {
			if exists(tried, twoDto1(r, c)) || board[r][c].editable {
				continue
			}
			tried = append(tried, twoDto1(r, c))
			board[r][c].editable = true
			if rand.Float32() < 0.5 {
				board[r][c].value = 0
			}
			count++
		}
		fmt.Println()
		fmt.Println(board.View(nextRandPos()))
		fmt.Println()
	})

}

// func TestWindow(t *testing.T) {
// 	lipgloss.SetColorProfile(termenv.TrueColor)
// 	t.Run("TestWindow", func(t *testing.T) {
// 		board := Board{}
// 		board.fill()
// 		row := lipgloss.JoinVertical(
// 			lipgloss.Right,
// 			"╭\n│",
// 			tab.Render("Lip Gloss"),
// 			activeTab.Render("BLush"),
// 			activeTab.Render("Eye Shadow"),
// 			tab.Render("Mascara"),
// 			tab.Render("Foundation"),
// 		)
// 		gap := tabGap.Render(strings.Repeat("\n", max(0, 20-lipgloss.Height(row)))) + "\n╰"
// 		row = lipgloss.JoinVertical(lipgloss.Bottom, row, gap)
// 		boardView := lipgloss.NewStyle().Height(lipgloss.Height(row) - 2).AlignVertical(lipgloss.Center).Render(board.View(0, 0))
// 		boardView = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderLeft(false).Render(boardView)
// 		row = lipgloss.JoinHorizontal(lipgloss.Top, row, boardView)
// 		fmt.Print(row + "\n\n")
// 	})
// }
