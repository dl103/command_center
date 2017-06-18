package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/brentnd/go-snowboy"
	"github.com/dl103/wav-player"
	"github.com/gordonklaus/portaudio"
)

var alertSound wavplayer.Player

func main() {
	// Setup Ctrl + C kill signal recognition
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	fmt.Println("Running recording. Ctrl-C to quit")

	// Setup audio stream
	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int16, 2048)
	inputStream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), &in)
	chk(err)
	defer inputStream.Close()
	chk(inputStream.Start())

	alertSound = wavplayer.NewPlayer("/Users/david/workspace/go_workspace/src/github.com/dl103/command_center/resources/beep.wav")

	// use snowboy to listen for hotword
	detector := SetupSnowboy()
	defer detector.Close()
	buf := new(bytes.Buffer)

	// Loop in stream
	for {
		chk(inputStream.Read())
		select {
		case <-sig:
			return
		default:
			binary.Write(buf, binary.LittleEndian, in)
			detector.ReadAndDetect(buf)
		}
	}
	chk(inputStream.Stop())
}

func handleDetection(result string) {
	fmt.Println("Detected", result)
	alertSound.Play()
	return
}

func SetupSnowboy() (d snowboy.Detector) {
	snowboyPath := "/Users/david/workspace/go_workspace/src/github.com/Kitt-AI/snowboy"
	resourceFile := snowboyPath + "/resources/common.res"
	modelFile := snowboyPath + "/resources/snowboy.umdl"

	d = snowboy.NewDetector(resourceFile)
	d.HandleFunc(snowboy.NewDefaultHotword(modelFile), handleDetection)
	d.HandleSilenceFunc(500*time.Millisecond, func(string) {
		fmt.Println("Silence detected")
	})
	return
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
