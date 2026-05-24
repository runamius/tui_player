package model

import (
	"audio_player/audio"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type Model struct {
	Choices   []string
	Cursor    int
	Current   string
	Playing   bool
	Ctrl      *beep.Ctrl
	Directory string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func InitialModel() Model {
	files, err := audio.GetMP3Files("/home/runamius/audio_player/music/")
	if err != nil {
		log.Fatal(err)
	}

	return Model{
		Choices:   files,
		Directory: "/home/runamius/audio_player/music/",
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter":
			if m.Ctrl == nil {
				selected := m.Choices[m.Cursor]
				ctrl := audio.PlayMusic(m.Directory + selected)
				m.Ctrl = ctrl
				m.Playing = true
			} else {
				speaker.Lock()
				m.Ctrl.Paused = !m.Ctrl.Paused
				speaker.Unlock()
				m.Playing = !m.Ctrl.Paused
			}

		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "Choose music:\n\n"
	for i, choice := range m.Choices {
		cursor := " "
		if i == m.Cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n"

	if m.Current != "" {
		s += "Current: " + m.Current
	}

	s += "\n\n↑ ↓ move | Enter play | q quit"

	return s
}
