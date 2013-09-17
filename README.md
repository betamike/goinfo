# goinfo

Package providing a proc-like interface to monitor running Go processes

# Dependencies

goinfo requires [go-fuse](https://github.com/hanwen/go-fuse) to implement the proc-like file system. Due to this dependency, goinfo currently only works on Linux.

To install `go-fuse` you can simply run
    $ go get github.com/hanwen/go-fuse

If you use [goat](https://github.com/mediocregopher/goat), a Goatfile is included in the project already.

# Usage

The project is not well documented right now, but now that I have an initial working version, I will be fleshing out the godocs soon. In the mean time:

To start simply import `goinfo` and start the info interface at a given mount point:

    package main

    import "github.com/betamike/goinfo"

    func main() {
      goinfo.Start("~/go/myapp")
      //App logic
    }

This will provide the following interface at `~/go/myapp`:

    |- myapp
      |- mem            // memory statistics (all found in runtime.MemStats)
        |- genmem       // general memory stats "<MemStats.Alloc> <MemStats.TotalAlloc> <MemStats.Sys> <Memstats.Lookups> <MemStats.Mallocs> <MemStats.Frees>" 
        |- heap         // heap memory stats "<MemStats.HeapAlloc> <MemStats.HeapSys> <MemStats.HeapIdle> <MemStats.HeapInuse> <MemStats.HeapReleased> <MemStats.HeapObjects>"
        |- stack        // stack memory stats "<MemStats.StackInuse> <MemStats.StackSys>"
        |- mspan        // mspan memory stats "<MemStats.MSpanInuse> <MemStats.MSpanSys>"
        |- mcache       // mcache memory stats "<MemStats.MCacheInuse> <MemStats.MCacheSys>"
        |- buckethash   // bucket hash info "<MemStats.BuckHashSys>"
        |- gc           // garbage collection stats  "<MemStats.NextGC> <MemStats.LastGC> <MemStats.PauseTotalNs> <MemStats.NumGC> <MemStats.EnableGC> <MemStats.DebugGC>"
      |- st
        |- stacktrace   // the current stacktrace of all goroutines (see runtime.Stack())

Check out the [runtime package](http://golang.org/pkg/runtime/) for more info.

Acessing this information is a simple as accessing a file:

    $ cat ~/go/myapp/mem/genmem
    817232 861568 271175600 65 541 146    

# License

This package is distributed under the MIT license.  See the LICENSE file for more details.
