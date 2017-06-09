package main

import (
  "fmt"
  "os"
  "os/signal"
  "github.com/gordonklaus/portaudio"
  "github.com/brentnd/go-snowboy"
)

func main() {
  // Setup Ctrl + C kill signal recognition
  sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

  fmt.Println("Running recording. Ctrl-C to quit")

  // Setup audio stream
  portaudio.Initialize()
  defer portaudio.Terminate()
  in := make([]int16, 64)
  stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	chk(err)
	defer stream.Close()
  chk(stream.Start())

  // use snowboy to listen for hotword
  detector := SetupSnowboy()
	defer detector.Close()

  // Loop in stream
  for {
    chk(stream.Read())
    select {
    case <-sig:
      return
    default:
    }
  }
  chk(stream.Stop())
}

func SetupSnowboy() (d snowboy.Detector) {
  // Take resource file as argument to NewDetector
  resourcePath := "~/workspace/go_workspace/src/github.com/Kitt-AI/snowboy/resources/common.res"
  d = snowboy.NewDetector(resourcePath)
  return
}

func chk(err error) {
  if err != nil {
    panic(err)
  }
}
