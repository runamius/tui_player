package model

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var defaultDir = `~`

func GetFiles(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("can't open dir, %s, %w", dirPath, err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	sort.Strings(dirs)

	var music []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".mp3") {
			music = append(music, entry.Name())
		}
	}
	sort.Strings(music)
	return append(dirs, music...), nil
}

func (m *Model) EnterDir(path string) error {
	path = strings.TrimSuffix(path, "/")

	filestat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if filestat.IsDir() {
		m.List, err = GetFiles(path)
		if err != nil {
			log.Printf("%v ", err)
			return nil
		}
		m.Directory = path + "/"
		m.Cursor = 0

		if len(m.List) > 0 {
			m.Current = m.List[0]
		} else {
			m.Current = ""
		}
	}
	return nil
}
