package gmfs

import (
	"os"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"github.com/betamike/goinfo/sources"
)

type GoMonitorFs struct {
	pathfs.FileSystem
	Sources map[string]sources.DataSource
}

func New() *GoMonitorFs {
	src := make(map[string]sources.DataSource)
	return &GoMonitorFs{pathfs.NewDefaultFileSystem(), src}
}

func (gfs *GoMonitorFs) AddSource(source sources.DataSource) {
	gfs.Sources[source.Name()] = source
}

func (gfs *GoMonitorFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	parts := strings.Split(name, string(os.PathSeparator))
	if len(parts) == 1 {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}
	if source, ok := gfs.Sources[parts[0]]; ok {
		size, updated, found := source.Metadata(parts[1])
		if !found {
			return nil, fuse.ENOENT
		}
		return &fuse.Attr{
			Mode:  fuse.S_IFREG | 0644,
			Size:  size,
			Mtime: updated,
		}, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (gfs *GoMonitorFs) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	if name == "" {
		listing := make([]fuse.DirEntry, 0, len(gfs.Sources))
		for name, _ := range gfs.Sources {
			listing = append(listing, fuse.DirEntry{Name: name, Mode: fuse.S_IFDIR})
		}
		return listing, fuse.OK
	}

	if source, ok := gfs.Sources[name]; ok {
		items := source.Listing()
		listing := make([]fuse.DirEntry, len(items))
		for i, item := range items {
			listing[i] = fuse.DirEntry{Name: item, Mode: fuse.S_IFREG}
		}
		return listing, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (gfs *GoMonitorFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}
	parts := strings.Split(name, string(os.PathSeparator))
	if source, ok := gfs.Sources[parts[0]]; ok {
		content, found := source.Contents(parts[1])
		if !found {
			return nil, fuse.EPERM
		}
		return nodefs.NewDataFile(content), fuse.OK
	}
	return nil, fuse.EPERM
}
