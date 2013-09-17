package sources

type DataSource interface {
    Name() string
    Listing() []string
    Contents(name string) ([]byte, bool)
    Metadata(name string) (uint64, uint64, bool)
}
