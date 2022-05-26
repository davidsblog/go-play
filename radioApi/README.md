# Radio API

A simple webservice written in go which plays radio streams when requested.

## Sample request for playing BBC Radio 1
```
POST http://localhost:8080/api/play

{ "url": "http://stream.live.vc.bbcmedia.co.uk/bbc_radio_one" }
```

## Sample request for stopping
```
POST http://localhost:8080/api/stop

{ }
```

## Sample request for volume up
```
POST http://localhost:8080/api/volume

{ "command": "+" }
```

## Mac and Linux
Tested for Mac and Linux (Ubuntu) and everything should build and run without any trouble.

On Ubuntu, I needed to install libasound2-dev:
`sudo apt install libasound2-dev`

## Windows:

Building on Windows requires gcc (for cgo) which can be 

### Install GCC with MinGW
https://code.visualstudio.com/docs/cpp/config-mingw 
