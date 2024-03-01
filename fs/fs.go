package fs

import (
	"context"
	"fmt"
	"strings"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"

	"github.com/hanwen/go-fuse/v2/fuse"
	pb "github.com/johnewart/freighter/freighter/proto"
	"zombiezen.com/go/log"
)

type FreighterRoot struct {
	fs.Inode
	counter    int
	Client     pb.FreighterClient
	Repository string
	Target     string
	RootPath   string
	Children   map[string]*fs.Inode
	Size       int64
	IsDir      bool
	Data       []byte
}

type FreighterNode struct {
	fs.Inode
	Name        string
	Client      pb.FreighterClient
	ContainerId string
	Repository  string
	Target      string
	Size        int64
	Data        []byte
}

func (r *FreighterRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = r.Mode()
	out.Size = uint64(r.Size)
	return 0
}

func (r *FreighterRoot) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	log.Infof(ctx, "Readdir %v", r.RootPath)
	if files, err := r.Client.GetDir(ctx, &pb.DirRequest{Repository: r.Repository, Target: r.Target, Path: r.RootPath}); err != nil {
		log.Errorf(ctx, "Error: %v", err)
		return nil, syscall.EIO
	} else {
		inodes := make([]*fs.Inode, 0)
		filenames := make([]string, 0)
		r.Children = make(map[string]*fs.Inode, len(files.Files))
		for i, f := range files.Files {
			filename := r.PathTo(f.Name)
			log.Infof(ctx, "Readdir: %s -> %s", f.Name, filename)
			var inode *fs.Inode
			if f.IsDir {
				inode = r.NewInode(ctx, &FreighterRoot{Client: r.Client, Repository: r.Repository, Target: r.Target, RootPath: r.PathTo(f.Name)}, fs.StableAttr{Mode: syscall.S_IFDIR, Ino: uint64(i)})
			} else {
				inode = r.NewInode(ctx, &FreighterNode{Client: r.Client, Repository: r.Repository, Target: r.Target, Name: f.Name, Size: f.Size, ContainerId: r.RootPath}, fs.StableAttr{Mode: 0755, Ino: uint64(i)})
			}
			r.Children[f.Name] = inode
			inodes = append(inodes, inode)
			filenames = append(filenames, f.Name)
			log.Infof(ctx, "Readdir inode: %s -> %v", f.Name, *inode)
		}

		return &FreighterDir{
			ctx:       ctx,
			root:      r,
			filenames: filenames,
		}, 0

	}
}

type FreighterDir struct {
	fs.DirStream
	ctx       context.Context
	root      *FreighterRoot
	filenames []string
}

func (d *FreighterDir) Next() (fuse.DirEntry, syscall.Errno) {
	fname := d.filenames[d.root.counter-1]
	inode := d.root.Children[fname]
	if inode != nil {
		log.Infof(d.ctx, "Next inode: %s -> %v", fname, *inode)
		return fuse.DirEntry{
			Mode: inode.Mode(),
			Name: fname,
			Ino:  inode.StableAttr().Ino,
		}, 0
	} else {
		log.Infof(d.ctx, "Next() no entry: %s", fname)
		return fuse.DirEntry{}, syscall.ENOENT
	}
}

func (d *FreighterDir) Close() {
	log.Infof(d.ctx, "Close: %v", d.root.counter)
	d.root.counter = 0
}

func (d *FreighterDir) HasNext() bool {
	log.Infof(d.ctx, "HasNext: %v", d.root.counter)
	if d.root.counter < len(d.filenames) {
		d.root.counter++
		return true
	}
	return false
}

func (r *FreighterRoot) PathTo(fname string) string {
	return fmt.Sprintf("%s/%s", r.RootPath, fname)
}

func (r *FreighterRoot) Mode() uint32 {
	if r.IsDir {
		return syscall.S_IFDIR | 0755
	} else {
		return 0755
	}
}

func (r *FreighterRoot) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Infof(ctx, "Lookup %v", r.PathTo(name))
	if i, ok := r.Children[name]; !ok {
		return nil, syscall.ENOENT
	} else {
		return i, 0
	}
}

func (r *FreighterRoot) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	_, file := strings.Split(r.RootPath, "/")[0], strings.Split(r.RootPath, "/")[1]
	log.Infof(ctx, "OPEN %s:%s /%s", r.Repository, r.Target, file)
	if r.IsDir {
		return nil, 0, syscall.EISDIR
	} else {
		if r.Data == nil {
			resp, err := r.Client.GetFile(ctx, &pb.FileRequest{Repository: r.Repository, Target: r.Target, Path: file})
			if err != nil {
				return nil, 0, syscall.EIO
			}
			r.Data = resp.Data
			log.Infof(ctx, "Read %d bytes", len(r.Data))
		}

		// We don't return a filehandle since we don't really need
		// one.  The file content is immutable, so hint the kernel to
		// cache the data.
		return nil, fuse.FOPEN_KEEP_CACHE, 0

	}
}

func (f *FreighterNode) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = f.Mode()
	out.Size = uint64(f.Size)
	return 0
}

func (f *FreighterNode) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	log.Infof(ctx, "OPEN %s:%s /%s", f.Repository, f.Target, f.Name)
	if f.Data == nil {
		resp, err := f.Client.GetFile(ctx, &pb.FileRequest{Repository: f.Repository, Target: f.Target, Path: f.Name})
		if err != nil {
			return nil, 0, syscall.EIO
		}
		f.Data = resp.Data
		log.Infof(ctx, "Read %d bytes", len(f.Data))
	}

	// We don't return a filehandle since we don't really need
	// one.  The file content is immutable, so hint the kernel to
	// cache the data.
	return nil, fuse.FOPEN_KEEP_CACHE, 0
}

func (f *FreighterNode) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	log.Infof(ctx, "Read %v@%d", f.Name, off)
	end := int(off) + len(dest)
	if end > len(f.Data) {
		end = len(f.Data)
	}
	return fuse.ReadResultData(f.Data[off:end]), 0
}

//var _ = (fs.NodeOnAdder)((*FreighterRoot)(nil))

var _ = (fs.NodeReaddirer)((*FreighterRoot)(nil))
var _ = (fs.NodeLookuper)((*FreighterRoot)(nil))
var _ = (fs.NodeGetattrer)((*FreighterRoot)(nil))
var _ = (fs.NodeOpener)((*FreighterRoot)(nil))

var _ = (fs.NodeGetattrer)((*FreighterNode)(nil))
var _ = (fs.NodeOpener)((*FreighterNode)(nil))
var _ = (fs.NodeReader)((*FreighterNode)(nil))
