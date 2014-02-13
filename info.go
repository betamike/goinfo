package goinfo

import (
	"net"
	"bufio"
	"strings"

	"github.com/betamike/goinfo/sources"
	"github.com/betamike/goinfo/sources/memstats"
	"github.com/betamike/goinfo/sources/stacktrace"
)

// Starts a goinfo endpoint at the given address. Also takes in any additional,
// custom data-sources that wish to be included in the endpoint. Returns a
// stop channel or an error. If close() is called on the stop channel, the
// endpoint will close and stop serving requests.
func Start(addr string, extraSrc ...sources.DataSource) (chan struct{}, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	stopCh := make(chan struct{})
	srcs := map[string]sources.DataSource{
		"stacktrace": stacktrace.NewStacktraceSource(),
		"memstats": memstats.NewMemStatsSource(),
	}
	for _, src := range extraSrc {
		srcs[src.Name()] = src
	}

	go func(){
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			// Spawn a routine to read from the new connection and handle its
			// request
			go handleConn(c, srcs)

			// If the stop channel is closed or written to, we close the whole
			// endpoint
			select {
			case <-stopCh:
				l.Close()
				return
			default:
			}
		}
	}()
	return stopCh, nil
}

func handleConn(c net.Conn, srcs map[string]sources.DataSource) {
	// Close no matter what happens in this function
	defer c.Close()

	// Read a request line
	buf := bufio.NewReader(c)
	request, err := buf.ReadString('\n')
	if err != nil {
		c.Close()
		return
	}

	// Split it on /. There will be a leading /, ignore that one
	requestTrimmed := strings.TrimRight(request, "\n")
	parts := strings.SplitN(requestTrimmed[1:], "/", 2)
	srcName := parts[0]
	srcContent := parts[1]

	// If the request is valid, write the data
	if src, ok := srcs[srcName]; ok {
		if data, ok := src.Contents(srcContent); ok {
			c.Write(data)
		}
	}
}
