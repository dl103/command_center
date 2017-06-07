package main

import (
  "fmt"
  "strconv"
  "github.com/HardWareGuy/portaudio-go"
)

func main() {
  fmt.Printf("Running recording\n")
  fmt.Printf("Version: " + strconv.Itoa(portaudio.Version()) + "\n")
}
