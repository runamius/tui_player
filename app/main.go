package main

import (
	"audio_player/model"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := model.InitialModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
