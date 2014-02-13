package stacktrace

import (
	"runtime"
	"time"
)

type StacktraceSource struct {
	getStacktraceRec chan chan StacktraceRecord
}

func NewStacktraceSource() *StacktraceSource {
	getCh := make(chan chan StacktraceRecord)
	source := &StacktraceSource{getCh}
	go source.loadStacktrace()
	return source
}

func (ss *StacktraceSource) Name() string {
	return "st"
}

func (ss *StacktraceSource) Listing() []string {
	return []string{"stacktrace"}
}

func (ss *StacktraceSource) Contents(name string) ([]byte, bool) {
	if name != "stacktrace" {
		return nil, false
	}
	stack, _ := ss.getStacktrace()
	return stack, true
}

func (ss *StacktraceSource) getStacktrace() ([]byte, time.Time) {
	ret := make(chan StacktraceRecord)
	ss.getStacktraceRec <- ret
	result := <-ret
	return result.Stacktrace, result.LastUpdated
}

func (ss *StacktraceSource) loadStacktrace() {
	var lastUpdated time.Time
	var size int = 512
	var i int = 0
	var stacktrace []byte
	for {
		select {
		case ret := <-ss.getStacktraceRec:
			now := time.Now()
			if now.After(lastUpdated.Add(2 * time.Second)) {
				for {
					stacktrace = make([]byte, size)
					i = runtime.Stack(stacktrace, true)
					if i == size {
						size = int(float64(size) * 1.5)
						continue
					}
					size = int(float64(i) * 1.25)
					break
				}
				lastUpdated = now
			}
			ret <- StacktraceRecord{stacktrace[:i], lastUpdated}
		}
	}
}
