package model

import (
	"audio_player/audio"
	"log"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faiface/beep/speaker"
)

// Player mode
func (m Model) Playermode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	case "ctrl+c":
		return m, tea.Quit

	case "down":
		m.moveDown()

	case "up":
		m.moveUp()

	case "enter":

		m.toggleTrack()

	case "+", "=":
		m.changeVolume(0.125)

	case "-", "_":
		m.changeVolume(-0.125)
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

	if newVolume < 0 {
		newVolume = 0
	}
	m.VolumeVar = math.Round(newVolume*10) / 10
	m.applyToVolume()
}

// Insert mode
