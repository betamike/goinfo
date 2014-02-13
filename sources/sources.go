package sources

// DataSource is an interface which provides information to be presented through
// goinfo. Client programs can implement their own DataSource types in order to
// expose additional information.
type DataSource interface {

	// Name returns a string identifier for this DataSource.  This will become
	// the name of the directory in the goinfo request.
	Name() string

	// Contents takes a string and returns a []byte representing the contents of
	// the data, and a bool indicating whether the provided name is actually
	// represented by this DataSource.
	Contents(name string) ([]byte, bool)
}
