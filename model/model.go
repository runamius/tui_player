package model

import (
	"audio_player/audio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	List      []string
	Cursor    int
	Current   string
	Playing   bool
	Directory string
	Player    *audio.Player
	VolumeVar float64

	InputMode bool
	Input     textinput.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func InitialModel() Model {
	ti := textinput.New()

	ti.Placeholder = "/home/user/Music"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	return Model{
		InputMode: true,
		Input:     ti,
	}
}

func (m *Model) loadDir() {
	dir := m.Input.Value()

	if dir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Cannot get home dir: %v", err)
			return
		}
		dir = home
	}
	if dir == "" {
		m.EnterDir("/home/user/Music")
	}

	err := m.EnterDir(dir)
	if err != nil {

	}
	m.InputMode = false
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		if m.InputMode {

			var cmd tea.Cmd
			m.Input, cmd = m.Input.Update(msg)

			switch msg.String() {

			case "ctrl+c":
				return m, tea.Quit

			case "enter":
				m.loadDir()
			}

			return m, cmd
		}

		return m.Playermode(msg)
	}

	return m, nil
}

func (m Model) View() string {

	if m.InputMode {
		return fmt.Sprintf(
			"Enter music directory:\n\n%s\n\nPress Enter",
			m.Input.View(),
		)
	}

	s := "Choose music:\n\n"

	for i, choice := range m.List {

		cursor := " "

		if i == m.Cursor {
			cursor = ">"
		}

		s += fmt.Sprintf(
			"%s %s\n",
			cursor,
			choice,
		)
	}

	s += "\n"

	if m.Current != "" {
		s += fmt.Sprintf("Current: %s   Cursor: %s", m.Current, m.List[m.Cursor])
	}

	s += "\n\n↑ ↓ move | Enter play | ctrl+c quit"
	bars := VolumeVisual(float64(m.VolumeVar))
	s += fmt.Sprintf("\nVolume: %s", bars)
	s += fmt.Sprintf("\n%s", m.Directory)

	return s
}

func VolumeVisual(Volume float64) string {
	size := 10
	bars := int(Volume + 5)
	halfbarCount := float64(Volume+5) - float64(bars)

	var halfbar string
	switch {
	case halfbarCount >= 0 && halfbarCount < 0.25:
		halfbar = " "
	case halfbarCount >= 0.25 && halfbarCount < 0.5:
		halfbar = "░"
	case halfbarCount >= 0.5 && halfbarCount < 0.75:
		halfbar = "▒"
	case halfbarCount >= 0.75 && halfbarCount < 1:
		halfbar = "▓"
	default:
		halfbar = ""
	}

	s := strings.Repeat("█", bars) + halfbar + strings.Repeat(" ", size-bars)
	return s
}
