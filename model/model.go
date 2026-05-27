package model

import (
	"audio_player/audio"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
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

	files, err := audio.GetMP3Files(dir)
	if err != nil {
		log.Printf("Cannot open directory: %v", err)
		return
	}

	if len(dir) > 0 && dir[len(dir)-1] != '/' {
		dir += "/"
	}

	m.Directory = dir
	m.List = files
	m.Cursor = 0
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

		// режим плеера
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

			if m.Player == nil || m.Player.Ctrl == nil {
				m.startTrack(selected)
			} else {
				if selected == m.Current {

					speaker.Lock()

					m.Playing = !m.Playing
					m.Player.Ctrl.Paused = !m.Playing

					speaker.Unlock()

				} else {
					speaker.Clear()
					m.startTrack(selected)
				}
			}

		case "+", "=":

			if m.VolumeVar < 5 {
				m.VolumeVar += 0.1
				m.VolumeVar = math.Round(m.VolumeVar*10) / 10
			}

			m.applyToVolume()

		case "-", "_":

			if m.VolumeVar > -5 {
				m.VolumeVar -= 0.125
				m.VolumeVar = math.Round(m.VolumeVar*10) / 10
			}

			m.applyToVolume()
		}
	}

	return m, nil
}

func (m *Model) startTrack(filename string) {
	path := m.Directory + filename
	p := audio.PlayMusic(path)
	if p != nil {
		m.Player = p
		m.Current = filename
		m.Playing = true
		m.applyToVolume()
	} else {
		log.Printf("Failed to play: %s", path)
	}
}

func (m *Model) applyToVolume() {
	if m.Player == nil || m.Player.Volume == nil {
		return
	}

	speaker.Lock()
	defer speaker.Unlock()

	m.Player.Volume.Volume = m.VolumeVar

	if m.VolumeVar <= -5.0 {
		m.Player.Volume.Silent = true
	} else {
		m.Player.Volume.Silent = false
	}
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
		s += "Current: " + m.Current
	}

	s += "\n\n↑ ↓ move | Enter play | ctrl+c quit"
	bars := VolumeVisual(float64(m.VolumeVar))
	s += fmt.Sprintf("\nVolume: %s", bars)

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
