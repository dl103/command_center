package main

import (
  "fmt"
  "os"
  // "strconv"
  // "github.com/HardWareGuy/portaudio-go"
  "github.com/brentnd/go-snowboy"
)

func main() {
  fmt.Printf("Running recording\n")
  // fmt.Printf("Version: " + strconv.Itoa(portaudio.Version()) + "\n")
  // portaudio.Initialize()
  // deviceInfo, _ := portaudio.DefaultInputDevice()
  // fmt.Printf("DeviceInfo: %+v\n", deviceInfo)
  // Look for resource file
  d := snowboy.NewDetector(os.Args[1])
	defer d.Close()

  // use snowboy to listen for hotword


  // beep
}
