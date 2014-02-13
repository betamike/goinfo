# goinfo

Package providing a tcp interface to monitor running Go processes

Check out the [docs][docs]. The [runtime][runtime] docs might be
helpful too.

# Usage

To start simply import `goinfo` and start the info interface at a given port:

    package main

    import "github.com/betamike/goinfo"

    func main() {
      goinfo.Start("localhost:10000")
      //App logic
    }

This will provide the following interface at port 10000:

```
|- memstats       // memory statistics (all found in runtime.MemStats)
  |- genmem       // general memory stats "<MemStats.Alloc> <MemStats.TotalAlloc> <MemStats.Sys> <Memstats.Lookups> <MemStats.Mallocs> <MemStats.Frees>"
  |- heap         // heap memory stats "<MemStats.HeapAlloc> <MemStats.HeapSys> <MemStats.HeapIdle> <MemStats.HeapInuse> <MemStats.HeapReleased> <MemStats.HeapObjects>"
  |- stack        // stack memory stats "<MemStats.StackInuse> <MemStats.StackSys>"
  |- mspan        // mspan memory stats "<MemStats.MSpanInuse> <MemStats.MSpanSys>"
  |- mcache       // mcache memory stats "<MemStats.MCacheInuse> <MemStats.MCacheSys>"
  |- buckethash   // bucket hash info "<MemStats.BuckHashSys>"
  |- gc           // garbage collection stats  "<MemStats.NextGC> <MemStats.LastGC> <MemStats.PauseTotalNs> <MemStats.NumGC> <MemStats.EnableGC> <MemStats.DebugGC>"
|- stacktrace
  |- stacktrace   // the current stacktrace of all goroutines (see runtime.Stack())
```

Acessing this information is a simple as writing/reading from a socket:

    $ echo /memstats/heap | nc localhost 10000
    595344 1048576 258048 790528 188416 326

# License

This package is distributed under the MIT license.  See the LICENSE file for more details.

[docs]: http://godoc.org/github.com/betamike/goinfo
[runtime]: http://golang.org/pkg/runtime/
