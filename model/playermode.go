package model

import (
	"audio_player/audio"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faiface/beep/speaker"
)

// Player mode
func (m Model) Playermode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {

	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyDown:
		m.moveDown()

	case tea.KeyUp:
		m.moveUp()

	case tea.KeyEnter:
		filestat, _ := os.Stat(m.Directory + m.List[m.Cursor])
		if filestat.IsDir() {
			m.EnterDir(m.Directory + m.List[m.Cursor])
		} else {
			m.toggleTrack()
		}

	case tea.KeyBackspace:
		m.GoToParentDir()

	case tea.KeyRunes:
		switch string(msg.Runes) {

		case "+", "=":
			m.changeVolume(0.125)

		case "-", "_":
			m.changeVolume(-0.125)
		case "m", "M", "Ь", "ь":
			m.Player.Volume.Silent = !m.Player.Volume.Silent

		case "1":
			m.InputMode = true
		}
	}
	return m, nil
}

func (m *Model) moveDown() {
	if m.Cursor < len(m.List)-1 {
		m.Cursor++
	}
}

func (m *Model) moveUp() {
	if m.Cursor > 0 {
		m.Cursor--
	}
}

func (m *Model) toggleTrack() {
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

func (m *Model) changeVolume(delta float64) {
	newVolume := m.VolumeVar + delta

	if newVolume > 5 {
		newVolume = 5
	}

	if newVolume < -5 {
		newVolume = -5
	}
	m.VolumeVar = math.Round(newVolume*10) / 10
	m.applyToVolume()
}

func (m *Model) GoToParentDir() {
	if m.Directory == "" {
		return
	}

	cleanPath := strings.TrimSuffix(m.Directory, "/")

	parent := filepath.Dir(cleanPath)

	if parent == cleanPath {
		return
	}

	err := m.EnterDir(parent)
	if err != nil {
		log.Printf("Failed to go to parent dir %s: %v", parent, err)
	}
}
