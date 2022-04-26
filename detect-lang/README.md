# detect-lang
A simple language detection microservice which could run as an Azure function.

## Test on Katacoda

Fire up a Katacoda Go environment, visit:

https://www.katacoda.com/scenario-examples/courses/environment-usages/go

Then do the following:

```
git clone https://github.com/davidsblog/go-play.git
cd go-play/detect-lang/
go install golang.org/dl/go1.18
go1.18 download
go1.18 build ./api/detect.go
./detect &
curl -X POST http://localhost:8080/api/language -d '{"Text":"hello, world"}'
```

## Useful references...

### Lingua (language detection)
https://github.com/pemistahl/lingua-go

### Azure functions with Go
https://www.thorsten-hans.com/azure-functions-with-go/

### Developing a RESTful API with Go and Gin
https://go.dev/doc/tutorial/web-service-gin