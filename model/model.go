package model

import (
	audio "audio_player"
	"fmt"
	"log"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faiface/beep/speaker"
)

type Model struct {
	List      []string
	Cursor    int
	Current   string
	Playing   bool
	Directory string
	Player    *audio.Player
	VolumeVar float64
}

func (m Model) Init() tea.Cmd {
	return nil
}

func InitialModel() Model {
	directory := "/home/runamius/audio_player/music/"
	files, err := audio.GetMP3Files(directory)
	if err != nil {
		log.Fatal(err)
	}

	return Model{
		List:      files,
		Directory: directory}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "down":
			if m.Cursor < len(m.List)-1 {
				m.Cursor++
			}
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "enter":
			selected := m.List[m.Cursor]
			if m.Player == nil || m.Player.Ctrl == nil || m.Player.Volume == nil {
				m.Player = audio.PlayMusic(m.Directory + m.List[m.Cursor])
				m.Playing = true
				m.Player.Volume.Volume = m.VolumeVar
			} else {
				if selected == m.Current {
					speaker.Lock()
					m.Playing = !m.Playing
					m.Player.Ctrl.Paused = !m.Playing
					speaker.Unlock()
				} else {
					speaker.Clear()
					m.Player = audio.PlayMusic(m.Directory + m.List[m.Cursor])
					m.Playing = true
					m.Player.Volume.Volume = m.VolumeVar
				}
			}
		case "+", "=":
			if m.VolumeVar < 5 {
				m.VolumeVar += 0.1
				m.VolumeVar = math.Round(m.VolumeVar*10) / 10
			}
			applyToVolume(&m)
		case "-", "_":
			if m.VolumeVar > -5 {
				m.VolumeVar -= 0.1
				m.VolumeVar = math.Round(m.VolumeVar*10) / 10
			}
			applyToVolume(&m)

		}

	}
	return m, nil
}

func applyToVolume(m *Model) {
	if m.Player != nil && m.Player.Ctrl != nil {
		speaker.Lock()
		if m.Player.Volume.Volume == -5.0 {
			m.Player.Volume.Silent = true
		} else {
			m.Player.Volume.Silent = false
		}
		m.Player.Volume.Volume = m.VolumeVar
		speaker.Unlock()
	}
}

func (m Model) View() string {
	s := "Choose music:\n\n"
	for i, choice := range m.List {
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

	s += "\n\n↑ ↓ move | Enter play | ctrl+c quit"
	s += fmt.Sprintf("\n volume: %.0f%%\n", m.VolumeVar*100/5)

	return s
}
