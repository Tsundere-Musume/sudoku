package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	easy = iota
	medium
	hard
	difficultyCount
)

const (
	difficulty = iota
	debug
	optionsCcount
)

type settings struct {
	Difficulty int  `json:"difficulty"`
	Debug      bool `json:"debug"`
}

type settingsModel struct {
	settings *settings
	cursor   int
}

func parseSettings(filepath string) (*settings, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	s := new(settings)
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *settingsModel) isValid() bool {
	if s.settings.Difficulty >= difficultyCount || s.settings.Difficulty < 0 {
		return false
	}
	return true
}
func (s *settingsModel) View() string {
	out := strings.Builder{}
	switch s.settings.Difficulty {
	case easy:
		out.WriteString("Difficulty: Easy")
	case medium:
		out.WriteString("Difficulty: Medium")
	case hard:
		out.WriteString("Difficulty: Hard")
	default:
		out.WriteString("Difficulty: Easy")
		s.settings.Difficulty = easy
	}

	out.WriteString("\n")

	out.WriteString(fmt.Sprintf("Debug: %v", s.settings.Debug))
	return out.String()
}

func updateSettings(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.currWindow != tabSettings {
			break
		}

		switch msg.String() {
		case "h", "left":
			if m.game.options.cursor == difficulty {
				if m.game.options.settings.Difficulty > 0 {
					m.game.options.settings.Difficulty--
				}
			} else {
				m.game.options.settings.Debug = !m.game.options.settings.Debug
			}
		case "l", "right":
			if m.game.options.cursor == difficulty {
				if m.game.options.settings.Difficulty < difficultyCount-1 {
					m.game.options.settings.Difficulty++
				}
			} else {
				m.game.options.settings.Debug = !m.game.options.settings.Debug
			}
		case "j", "down":
			if m.game.options.cursor < optionsCcount-1 {
				m.game.options.cursor++
			}
		case "k", "up":
			if m.game.options.cursor > 0 {
				m.game.options.cursor--
			}
		case "enter":
			m.game.options.save()
		}
	}
	return m, nil
}

func (m *settingsModel) save() {
	file, _ := os.OpenFile("settings.json", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)
	_ = encoder.Encode(m.settings)
}
