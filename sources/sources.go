package sources

// DataSource is an interface which provides information to be
// presented through the goinfo file system. Client programs
// can implement their own DataSource types in order to expose
// additional information through the file system.
type DataSource interface {

	// Name returns a string identifier for this DataSource.
	// This will become the name of the directory in the
	// goinfo file system.
	Name() string

	// Listing returns a list of strings, each entry
	// representing the name of a file which will be
	// accessible from inside this DataSource's
	// directory.
	Listing() []string

	// Contents takes a string which will be one of the
	// the entries from the slice provided in Listings().
	// It returns a []byte representing the contents of
	// the file, and a bool indicating whether the
	// provided name is actually represented by this
	// DataSource.
	Contents(name string) ([]byte, bool)

	// Metadata takes a string which will be one of the
	// the entries from the slice provided in Listings().
	// It returns the expected file size, the last modified
	// timestamp as a Unix timestamp, and a bool indicating
	// whether the provided name is actually represented by
	// this DataSource.
	Metadata(name string) (uint64, uint64, bool)
}
