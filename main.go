package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
	openai "github.com/sashabaranov/go-openai"
	"os"
	"os/signal"
	"time"
)

const (
	sampleRate = 44100
	seconds    = 2
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing required argument:  output file name")
		return
	}
	fmt.Println("Recording.  Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	fileName := os.Args[1]
	f, err := os.Create(fileName)
	chk(err)

	defer func() {
		chk(f.Close())
	}()

	portaudio.Initialize()
	time.Sleep(1)
	defer portaudio.Terminate()
	in := make([]int16, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
loop:
	for {
		chk(stream.Read())
		chk(binary.Write(f, binary.LittleEndian, in))
		select {
		case <-sig:
			break loop
		default:
		}
	}
	chk(stream.Stop())

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateTranscription(
		context.Background(),
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: fileName,
		},
	)

	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return
	}
	fmt.Println(resp.Text)
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
