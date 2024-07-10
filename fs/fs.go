package fs

import (
	"context"
	"fmt"
	"hash/fnv"
	"path/filepath"
	"strings"
	"syscall"

	"google.golang.org/grpc"
	"zombiezen.com/go/log"

	"github.com/hanwen/go-fuse/v2/fs"

	"github.com/hanwen/go-fuse/v2/fuse"
	pb "github.com/johnewart/freighter/freighter/proto"
)

type FreighterRoot struct {
	fs.Inode
	counter    int
	Client     pb.FreighterClient
	Repository string
	Target     string
	Path       string
	Children   map[string]*fs.Inode
	Size       int64
	IsDir      bool
	Data       []byte
}

type FreighterNode struct {
	fs.Inode
	Name       string
	Client     pb.FreighterClient
	Path       string
	Repository string
	Target     string
	Size       uint64
	Data       []byte
	Mode       uint32
	Mtime      uint64
	Atime      uint64
	Ctime      uint64
	ExtraData  string
	Type       string
}

func (r *FreighterRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = r.Mode()
	out.Size = uint64(r.Size)
	return 0
}

type NamedInode struct {
	inode *fs.Inode
	name  string
}

func (r *FreighterRoot) PathTo(fname string) string {
	return fmt.Sprintf("%s/%s", r.Path, fname)
}

func (r *FreighterRoot) OnAdd(ctx context.Context) {
	log.Debugf(ctx, "OnAdd: %v", r.Path)
	if response, err := r.Client.GetTree(ctx, &pb.TreeRequest{Repository: r.Repository, Target: r.Target}); err != nil {
		log.Errorf(ctx, "Error: %v", err)
	} else {

		for _, file := range response.Files {
			dir, base := filepath.Split(filepath.Clean(file.Name))
			if base == "" {
				continue
			}

			p := r.EmbeddedInode()
			for _, comp := range strings.Split(dir, "/") {
				if len(comp) == 0 {
					continue
				}
				ch := p.GetChild(comp)
				if ch == nil {
					ch = p.NewPersistentInode(ctx,
						&fs.Inode{},
						fs.StableAttr{Mode: syscall.S_IFDIR})
					p.AddChild(comp, ch, false)
				}
				p = ch
			}

			fullPath := r.PathTo(file.Name)
			if fullPath[:2] == "//" {
				fullPath = fullPath[1:]
			}

			attr := fuse.Attr{
				Mode:  file.Mode,
				Size:  uint64(file.Size),
				Mtime: file.ModTime,
				Atime: file.AccessTime,
				Ctime: file.ChangeTime,
			}

			log.Debugf(ctx, "Adding {%s}", file.String())

			if file.Type == pb.FileType_SYMLINK {
				node := &fs.MemSymlink{
					Data: []byte(file.ExtraData),
					Attr: attr,
				}
				p.AddChild(base, r.NewPersistentInode(ctx, node, fs.StableAttr{Mode: file.Mode, Ino: 0}), false)
			} else {
				var fileType string
				switch file.Type {
				case pb.FileType_FILE:
					fileType = "F"
				case pb.FileType_DIR:
					fileType = "D"
				case pb.FileType_SYMLINK:
					fileType = "S"
				}

				node := &FreighterNode{
					Client:     r.Client,
					Repository: r.Repository,
					Target:     r.Target,
					Name:       file.Name,
					Size:       file.Size,
					Path:       fullPath,
					Mode:       file.Mode,
					Mtime:      file.ModTime,
					Atime:      file.AccessTime,
					Ctime:      file.ChangeTime,
					ExtraData:  file.ExtraData,
					Type:       fileType,
				}
				inode := r.NewPersistentInode(ctx, node, fs.StableAttr{Mode: file.Mode, Ino: 0})
				p.AddChild(base, inode, false)
			}
		}
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (f *FreighterNode) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = f.Mode
	out.Size = uint64(f.Size)
	out.Atime = f.Atime
	out.Mtime = f.Mtime
	out.Ctime = f.Ctime
	return 0
}

func (f *FreighterNode) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	log.Debugf(ctx, "OPEN %s:%s %s", f.Repository, f.Target, f.Path)
	if f.Data == nil {
		log.Debugf(ctx, "Fetching file %s from %s:%s", f.Path, f.Repository, f.Target)
		maxSizeOption := grpc.MaxCallRecvMsgSize(1024 * 1024 * 1024)
		resp, err := f.Client.GetFile(ctx, &pb.FileRequest{Repository: f.Repository, Target: f.Target, Path: f.Path}, maxSizeOption)
		if err != nil {
			log.Errorf(ctx, "Error fetching file data for %s:%s %s: %v", f.Repository, f.Target, f.Path, err)
			return nil, 0, syscall.EIO
		}
		f.Data = resp.Data
		log.Debugf(ctx, "Read %d bytes", len(f.Data))
	}

	// We don't return a filehandle since we don't really need
	// one.  The file content is immutable, so hint the kernel to
	// cache the data.
	return nil, fuse.FOPEN_KEEP_CACHE, 0
}

func (f *FreighterNode) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	log.Debugf(ctx, "Reading %v", f.Name)
	log.Debugf(ctx, "File is %d bytes", len(f.Data))
	log.Debugf(ctx, "Reading %d bytes from %d", len(dest), off)
	end := int(off) + len(dest)
	if end > len(f.Data) {
		end = len(f.Data)
	}
	return fuse.ReadResultData(f.Data[off:end]), 0
}

var _ = (fs.NodeOnAdder)((*FreighterRoot)(nil))

var _ = (fs.NodeGetattrer)((*FreighterNode)(nil))
var _ = (fs.NodeOpener)((*FreighterNode)(nil))
var _ = (fs.NodeReader)((*FreighterNode)(nil))
