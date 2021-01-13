package main

import (
	"os"
	"os/exec"
)

func main() {
	for {

		cmd := exec.Command("./webrtc-desktop-pion-offer.exe")
		cmd.Stdout = os.Stdout
		cmd.Run()
		cmd.Wait()
	}
}
