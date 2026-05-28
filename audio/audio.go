package audio

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var initSpeaker sync.Once

type Player struct {
	Ctrl   *beep.Ctrl
	Volume *effects.Volume
}

func PlayMusic(path string) *Player {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Printf("Error decoding mp3: %v\n", err)
		f.Close()
		return nil
	}

	initSpeaker.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	})

	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   1.0,
		Silent:   false,
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
