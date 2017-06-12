package main

import (
  "bytes"
  "fmt"
  "os"
  "os/signal"
  "time"

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
  in := make([]byte, 2048)
  reader := bytes.NewReader(in)
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
      fmt.Println(detector.ReadAndDetect(reader))
    }
  }
  chk(stream.Stop())
}

func handleDetection(result string) {
  fmt.Println("Detected", result)
  return
}

func SetupSnowboy() (d snowboy.Detector) {
  snowboyPath := "/Users/david/workspace/go_workspace/src/github.com/Kitt-AI/snowboy"
  resourceFile := snowboyPath + "/resources/common.res"
  modelFile := snowboyPath + "/resources/snowboy.umdl"

  d = snowboy.NewDetector(resourceFile)
	d.HandleFunc(snowboy.NewDefaultHotword(modelFile), handleDetection)
	d.HandleSilenceFunc(500*time.Millisecond, func(string) {
		fmt.Println("silence detected")
	})
  return
}

func chk(err error) {
  if err != nil {
    panic(err)
  }
}
