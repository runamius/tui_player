package audio

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	Ctrl   *beep.Ctrl
	Volume *effects.Volume
}

var initSpeaker sync.Once

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

func PlayMusic(file string) *Player {

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

	initSpeaker.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	})
	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2.0,
		Volume:   0.0,
	}
	ctrl := &beep.Ctrl{
		Streamer: volume,
		Paused:   false,
	}

	speaker.Play(ctrl)
	return &Player{
		Ctrl:   ctrl,
		Volume: volume,
	}

}
