package main

import (
	"io"
	"os/exec"
)

// example command:
// wget -O - http://stream.live.vc.bbcmedia.co.uk/bbc_radio_one | mpv -

// here are some streams:
// https://garfnet.org.uk/cms/tables/radio-frequencies/internet-radio-player/

func main() {
	wget := exec.Command(
		"wget", "-O", "-", "http://stream.live.vc.bbcmedia.co.uk/bbc_radio_one",
	)
	mpv := exec.Command("mpv", "-")
	mpv.Stdin, wget.Stdout = io.Pipe()
	wget.Start()
	mpv.Run()
}
