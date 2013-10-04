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

func (ds *CountsDataSource) Listing() []string {
	return []string{"countA", "countB"}
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

func (ds *CountsDataSource) Metadata(name string) (l uint64, ts uint64, found bool) {
	content, _ := ds.Contents(name)
	switch name {
	case "countA":
		found = true
		ts = ds.timestampA
		l = uint64(len(content))
	case "countB":
		found = true
		ts = ds.timestampB
		l = uint64(len(content))
	default:
		found = false
		ts = 0
		l = 0
	}
	return
}

func main() {
	countDS := &CountsDataSource{}

	goinfo.AddSource(countDS)
	err := goinfo.Start("example")
	if err != nil {
		panic("Oh no! could not mount our monitoring file system: " + err.Error())
	}
	defer goinfo.StopAll()

	//Start app business logic
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
