package main

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hajimehoshi/oto"
	"github.com/itchyny/volume-go"
	"github.com/tosone/minimp3"
)

var dec *minimp3.Decoder // global mp3 decoder
var stop chan bool       // for sending stop signals

type PlayRequest struct {
	Url string `json:"url"`
}

type VolumeRequest struct {
	Command string `json:"command"`
}

func playStream(url string) *minimp3.Decoder {
	var err error
	var response *http.Response

	if response, err = http.Get(url); err != nil {
		log.Fatal(err)
		return nil
	}

	var dec *minimp3.Decoder
	if dec, err = minimp3.NewDecoder(response.Body); err != nil {
		log.Fatal(err)
		return nil
	}
	<-dec.Started()

	log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)

	var waitForPlayOver = new(sync.WaitGroup)
	waitForPlayOver.Add(1)

	go func() {
		defer func() {
			dec.Close()
			dec = nil
			response.Body.Close()
		}()

		var context *oto.Context
		if context, err = oto.NewContext(dec.SampleRate, dec.Channels, 2, 4096); err != nil {
			log.Fatal(err)
			return
		}
		defer context.Close()
		var player = context.NewPlayer()
		defer player.Close()
		for {
			var data = make([]byte, 512)
			_, err = dec.Read(data)
			if len(stop) > 0 {
				log.Println("Stop signal waiting")
				break
			}
			if err == io.EOF {
				log.Println("EOF")
				break
			}
			if err != nil {
				log.Fatal(err)
				break
			}
			player.Write(data)
		}
		log.Println("Streaming stopped")
		waitForPlayOver.Done()

		if len(stop) > 0 {
			<-stop
			log.Println("Stop signal consumed")
		} else {
			go playStream(url) // try to restart a faulted stream
		}
	}()

	return dec
}

func postVol(c *gin.Context) {
	var request VolumeRequest

	if err := c.BindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "Bad request, could not get volume command")
		return
	}

	vol, err := volume.GetVolume()
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not get volume")
		return
	}

	if request.Command == "+" {
		vol++
	}
	if request.Command == "-" {
		vol--
	}
	err = volume.SetVolume(vol)
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not set volume")
		return
	}
	c.String(http.StatusOK, string(vol))
}

func postPlay(c *gin.Context) {
	var request PlayRequest

	if err := c.BindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "Bad request, could not get URL")
		return
	}

	if dec = playStream(request.Url); dec != nil {
		c.String(http.StatusOK, "OK")
		return
	}

	c.String(http.StatusInternalServerError, "Could not play radio")
}

func postStop(c *gin.Context) {
	if dec == nil {
		c.String(http.StatusBadRequest, "Nothing to stop")
		return
	}
	if len(stop) == 0 {
		stop <- true // signal a stop
	}
	c.String(http.StatusOK, "OK")
}

func main() {
	stop = make(chan bool, 1)

	engine := gin.Default()
	engine.POST("/api/play", postPlay)
	engine.POST("/api/stop", postStop)
	engine.POST("/api/volume", postVol)
	engine.Run(":8080")
}
