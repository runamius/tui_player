package audio

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func GetMP3Files(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("can't open dir, %s, %w", dirPath, err)
	}

	var files []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToLower(entry.Name()), ".mp3") {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)
	return files, nil
}

func PlayMusic(file string) *beep.Ctrl {

	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return nil
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Println(err)
		return nil
	}

	ctrl := &beep.Ctrl{
		Streamer: streamer,
		Paused:   false,
	}

	speaker.Init(
		format.SampleRate,
		format.SampleRate.N(time.Second/10),
	)

	speaker.Play(ctrl)

	return ctrl
}
