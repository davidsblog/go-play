package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gordonklaus/portaudio"
)

// see:
// https://github.com/valentijnnieman/audio_streamer/blob/master/audio_streamer_client/main.go

// streams:
// https://garfnet.org.uk/cms/tables/radio-frequencies/internet-radio-player/

const sampleRate = 44100
const seconds = 2

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	buffer := make([]float32, sampleRate*seconds)

	stream, err := portaudio.OpenDefaultStream(0, 1, sampleRate, len(buffer), func(out []float32) {
		resp, err := http.Get("http://stream.live.vc.bbcmedia.co.uk/bbc_radio_one")
		chk(err)
		body, _ := ioutil.ReadAll(resp.Body)
		responseReader := bytes.NewReader(body)
		binary.Read(responseReader, binary.BigEndian, &buffer)
		for i := range out {
			out[i] = buffer[i]
		}
	})
	chk(err)
	chk(stream.Start())
	time.Sleep(time.Second * 40)
	chk(stream.Stop())
	defer stream.Close()

	if err != nil {
		fmt.Println(err)
	}

}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
