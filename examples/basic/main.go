package main

import (
	"github.com/betamike/goinfo"
)

func main() {
	_, err := goinfo.Start("localhost:10000")
	if err != nil {
		panic("Oh no! could not mount our endpoint: " + err.Error())
	}

	// Start app business logic
	// Use the endpoint with:
	// echo /stacktrace/stacktrace | nc localhost 10000
	select {}
}
