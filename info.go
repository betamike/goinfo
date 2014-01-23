package goinfo

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"github.com/betamike/goinfo/gmfs"
	"github.com/betamike/goinfo/sources"
	"github.com/betamike/goinfo/sources/memstats"
	"github.com/betamike/goinfo/sources/stacktrace"
)

var gfs *gmfs.GoMonitorFs
var servers map[string]*fuse.Server

func init() {
	servers = make(map[string]*fuse.Server)

	//unmount all when program exits cleanly or uncleanly
	sig := make(chan os.Signal, 1)
	go func() {
		<-sig
		StopAll()
		os.Exit(1)
	}()
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	//create the file system
	gfs = gmfs.New()
	gfs.AddSource(stacktrace.NewStacktraceSource())
	gfs.AddSource(memstats.NewMemStatsSource())
}

// AddSource takes a DataSource and adds it to the goinfo file system.
// This allows programs to define additional types of information to be
// made accessible via the goinfo file system.
func AddSource(source sources.DataSource) {
	gfs.AddSource(source)
}

// Start takes a path to the directory where the goinfo file system
// should be mounted. If the directory does not exist, it will be 
// created. Start will return an error if the directory cannot be
// created or if the file system cannot be mounted at this location
// for any reason.
func Start(mountpoint string) error {
	//already mounted there
	if _, found := servers[mountpoint]; found {
		return nil
	}

	if _, err := os.Stat(mountpoint); os.IsNotExist(err) {
		if err = os.Mkdir(mountpoint, 0755); err != nil {
			return err
		}
	}

	nfs := pathfs.NewPathNodeFs(gfs, nil)
	conn := nodefs.NewFileSystemConnector(nfs.Root(), nil)
	server, err := fuse.NewServer(conn.RawFS(), mountpoint, &fuse.MountOptions{AllowOther: true})
	if err != nil {
		return errors.New("Failed to mount monitoring fs at " + mountpoint + ": " + err.Error())
	}

	servers[mountpoint] = server

	//start handling the fs calls
	go server.Serve()

	return nil
}

// Stop will unmount the goinfo file system mounted at the provided
// path. If the path is not a goinfo mounted file system or if the 
// unmounted fails, Stop will return an error.
func Stop(mountpoint string) error {
	server, found := servers[mountpoint]
	if !found {
		return errors.New("No file system mounted at " + mountpoint)
	}

	err := server.Unmount()
	if err != nil {
		return err
	}

	delete(servers, mountpoint)
	return nil
}

// StopAll will unmount all goinfo file systems mounted by this
// program.
func StopAll() {
	for _, server := range servers {
		server.Unmount()
	}
	servers = make(map[string]*fuse.Server)
}
