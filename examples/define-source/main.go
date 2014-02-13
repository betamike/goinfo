package main

import (
	"fmt"
	"github.com/betamike/goinfo"
	"time"
)

type CountsDataSource struct {
	arbitraryCountA int
	timestampA      uint64

	arbitraryCountB int
	timestampB      uint64
}

func (ds *CountsDataSource) Name() string {
	return "counts"
}

func (ds *CountsDataSource) Contents(name string) (content []byte, found bool) {
	switch name {
	case "countA":
		found = true
		content = []byte(fmt.Sprintf("%d\n", ds.arbitraryCountA))
	case "countB":
		found = true
		content = []byte(fmt.Sprintf("%d\n", ds.arbitraryCountB))
	default:
		found = false
		content = nil
	}
	return
}

// Use the endpoint with:
// echo /counts/countA | nc localhost 10000
func main() {
	countDS := &CountsDataSource{}

	_, err := goinfo.Start("localhost:10000", countDS)
	if err != nil {
		panic("Oh no! could not mount our endpoint: " + err.Error())
	}

	// Start app business logic
	for {
		select {
		case <-time.After(500 * time.Millisecond):
			ts := uint64(time.Now().Unix())
			if ts%3 == 0 {
				countDS.arbitraryCountA += 1
				countDS.timestampA = ts
			} else {
				countDS.arbitraryCountB += 1
				countDS.timestampB = ts
			}
		}
	}
}
