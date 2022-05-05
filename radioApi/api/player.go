package main

import (
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/itchyny/volume-go"
)

// example command:
// wget -O - http://stream.live.vc.bbcmedia.co.uk/bbc_radio_one | mpv -

// here are some streams:
// https://garfnet.org.uk/cms/tables/radio-frequencies/internet-radio-player/

func play(url string, ch chan os.Process) {
	wget := exec.Command("wget", "-O", "-", url)
	mpv := exec.Command("mpv", "-")
	mpv.Stdin, wget.Stdout = io.Pipe()
	wget.Start()
	mpv.Start()
	ch <- *wget.Process
	ch <- *mpv.Process
}

func main() {
	vol, err := volume.GetVolume()
	if err != nil {
		panic("cannot get volume info")
	}
	err = volume.SetVolume(vol / 2)
	if err != nil {
		panic("cannot set volume")
	}

	ch := make(chan os.Process)
	go play("http://stream.live.vc.bbcmedia.co.uk/bbc_radio_one", ch)
	wget := <-ch
	mpv := <-ch

	time.Sleep(time.Second * 30)

	mpv.Kill()
	wget.Kill()
}
