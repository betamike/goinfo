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

	The above DataSource would add this structure to the goinfo file system:
	The following directory structure is found at the mounted path:
		|- mount
		    |- counts
		        |- countA
		        |- countB

	For the full example check out "examples/define-source" in the source repository.
*/
package goinfo
