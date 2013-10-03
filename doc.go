/*
	Package goinfo provides a proc-like interface to inspect a running Go program.
	by default a memory and stacktrace source are enabled.

	The goinfo file system can be easily mounted with the defaults:
		goinfo.Start("/path/to/mount")

	The following directory structure is found at the mounted path:
		|- mount
		    |- mem              // memory statistics (all found in runtime.MemStats)
		        |- genmem       // general memory stats "<MemStats.Alloc> <MemStats.TotalAlloc> <MemStats.Sys> <Memstats.Lookups> <MemStats.Mallocs> <MemStats.Frees>"
		        |- heap         // heap memory stats "<MemStats.HeapAlloc> <MemStats.HeapSys> <MemStats.HeapIdle> <MemStats.HeapInuse> <MemStats.HeapReleased> <MemStats.HeapObjects>"
		        |- stack        // stack memory stats "<MemStats.StackInuse> <MemStats.StackSys>"
		        |- mspan        // mspan memory stats "<MemStats.MSpanInuse> <MemStats.MSpanSys>"
		        |- mcache       // mcache memory stats "<MemStats.MCacheInuse> <MemStats.MCacheSys>"
		        |- buckethash   // bucket hash info "<MemStats.BuckHashSys>"
		        |- gc           // garbage collection stats "<MemStats.NextGC> <MemStats.LastGC> <MemStats.PauseTotalNs> <MemStats.NumGC> <MemStats.EnableGC> <MemStats.DebugGC>"
		    |- st
		        |- stacktrace   // the current stacktrace of all goroutines (see runtime.Stack())

	Client programs can implement the DataSource interface, which will allow them to make addiitonal
	information accessible via this file system. Each DataSource will be represented as a directory
	under the root mount point, and can provide any number of files beneath them. For example:

		type FooDataSource struct {
			numBars int
			tsBar uint64
			numBats int
			tsBat uint64
		}
		
		func (f *FooDataSource) Name() string {
			return "foo"
		}
		
		func (f *FooDataSource) Listing() []string {
			return []string{"bar", "bat"}
		}
		
		func (f *FooDataSource) Contents(name string) (content []byte, found bool) {
			switch (name) {
				case "bar":
					found = true
					content = []byte(fmt.Printf("%d", f.numBars))
				case "bat":
					found = true
					content = []byte(fmt.Printf("%d", f.numBats))
				default:
					found = false
					content = nil
			}
			return
		}
		
		func (f *FooDataSource) Name(name string) (l uint64,  ts uint64, found bool) {
			switch (name) {
				case "bar":
					found = true
					ts = f.tsBar
					l = uint64(len(fmt.Printf("%d", f.numBars)))
				case "bat":
					found = true
					ts = f.tsBat
					l  = uint64(len(fmt.Printf("%d", f.numBars)))
				default:
					found = false
					ts = 0
					l = 0
			}
			return
		}

	The above DataSource would add this structure to the goinfo file system:
	The following directory structure is found at the mounted path:
		|- mount
		    |- foo
		        |- bar
		        |- bat
*/
package goinfo
