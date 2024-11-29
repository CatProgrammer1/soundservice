package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Sound struct {
	file    *os.File
	control *beep.Ctrl
}

func new(path string) *Sound {
	sound := &Sound{}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	sound.file = file

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		file.Close()
		panic(err)
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	if err != nil {
		file.Close()
		panic(err)
	}

	sound.control = &beep.Ctrl{Streamer: streamer, Paused: false}

	return sound
}

func (sound *Sound) Play() {
	sound.Unpause()

	speaker.Play(sound.control)
}

func (sound *Sound) Stop() {
	sound.Pause()
	sound.control.Streamer.(beep.StreamSeekCloser).Seek(0)
}

func (sound *Sound) Pause() {
	sound.control.Paused = true
}

func (sound *Sound) Destroy() {
	sound.file.Close()
	sound.control.Streamer.(beep.StreamSeekCloser).Close()
}

func (sound *Sound) Unpause() {
	sound.control.Paused = false
}
