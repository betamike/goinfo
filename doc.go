/*
	Package goinfo provides a tcp interface to inspect a running Go program. By
	default a memory and stacktrace source are enabled.

	The goinfo endpoint can be easily setup with the defaults:
		goinfo.Start("localhost:10000").

	Requests are made by first writing a newline delimited path to the endpoint,
	then reading all data from the endpoint. The endpoint will close the
	connection when it is done writing. An example request would be:

	"/stacktrace/stacktrace\n"
	or
	"/memstats/gc\n"

	The "directory" is the name of the data source, and the "file" is what part
	of that source you want to read.

	The following data sources/parts are available:
		|- memstats         // memory statistics (all found in runtime.MemStats)
			|- genmem       // general memory stats "<MemStats.Alloc> <MemStats.TotalAlloc> <MemStats.Sys> <Memstats.Lookups> <MemStats.Mallocs> <MemStats.Frees>"
			|- heap         // heap memory stats "<MemStats.HeapAlloc> <MemStats.HeapSys> <MemStats.HeapIdle> <MemStats.HeapInuse> <MemStats.HeapReleased> <MemStats.HeapObjects>"
			|- stack        // stack memory stats "<MemStats.StackInuse> <MemStats.StackSys>"
			|- mspan        // mspan memory stats "<MemStats.MSpanInuse> <MemStats.MSpanSys>"
			|- mcache       // mcache memory stats "<MemStats.MCacheInuse> <MemStats.MCacheSys>"
			|- buckethash   // bucket hash info "<MemStats.BuckHashSys>"
			|- gc           // garbage collection stats "<MemStats.NextGC> <MemStats.LastGC> <MemStats.PauseTotalNs> <MemStats.NumGC> <MemStats.EnableGC> <MemStats.DebugGC>"
		|- stacktrace
			|- stacktrace   // the current stacktrace of all goroutines (see runtime.Stack())

	Client programs can implement the DataSource interface, which will allow
	them to make addiitonal information accessible via this endpoint. Each
	DataSource will be represented as a directory, and can provide any number of
	data parts.

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

	The above DataSource would add these directory/parts to the endpoint:
		|- counts
			|- countA
			|- countB

	For the full example check out "examples/define-source" in the source repository.
*/
package goinfo
