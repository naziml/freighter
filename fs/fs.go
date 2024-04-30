package fs

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"syscall"
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
	Size       int64
	Data       []byte
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

// (f *FreighterRoot) OnAdd(ctx context.Context) {
//}

func (f *FreighterRoot) LoadTree(ctx context.Context) error {
	if response, err := f.Client.GetTree(ctx, &pb.TreeRequest{Repository: f.Repository, Target: f.Target}); err != nil {
		log.Errorf(ctx, "Error: %v", err)
		return err
	} else {
		nodeTree := make(map[string]*NamedInode, 0)
		treeLevels := make(map[int][]*NamedInode, 0)
		maxLevel := 0

		for _, file := range response.Files {
			log.Infof(ctx, "Registering: %v", file)
			pathParts := strings.Split(file.Name, "/")
			inode := NamedInode{
				inode: f.NewInode(ctx, &FreighterRoot{Client: f.Client, IsDir: true, Repository: f.Repository, Target: f.Target, Path: f.PathTo(file.Name)}, fs.StableAttr{Mode: syscall.S_IFDIR, Ino: uint64(file.Size)}),
				name:  file.Name,
			}
			nodeTree[file.Name] = &inode
			treeLevels[len(pathParts)] = append(treeLevels[len(pathParts)], &inode)
			maxLevel = max(maxLevel, len(pathParts))
		}

		for i := 0; i < maxLevel; i++ {
			if i == 0 {
				for _, inode := range treeLevels[i] {
					f.AddChild(inode.name, inode.inode, false)
				}
			} else {
				for _, inode := range treeLevels[i] {
					parentPath := strings.Join(strings.Split(inode.name, "/")[:i], "/")
					parent := nodeTree[parentPath]
					parent.inode.AddChild(inode.name, inode.inode, false)
				}
			}
		}
	}

	return nil

}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (r *FreighterRoot) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	log.Infof(ctx, "Readdir: Reading listing for %v", r.Path)
	if files, err := r.Client.GetDir(ctx, &pb.DirRequest{Repository: r.Repository, Target: r.Target, Path: r.Path}); err != nil {
		log.Errorf(ctx, "Error: %v", err)
		return nil, syscall.EIO
	} else {
		inodes := make([]*fs.Inode, 0)
		filenames := make([]string, 0)
		r.Children = make(map[string]*fs.Inode, len(files.Files))
		for _, f := range files.Files {
			var inode *fs.Inode
			f.Name = strings.TrimSuffix(f.Name, "/")
			f.Name = strings.Replace(f.Name, r.Path, "", 1)
			f.Name = strings.TrimPrefix(f.Name, "/")
			fullPath := r.PathTo(f.Name)
			inodeId := uint64(hash(fullPath))

			log.Infof(ctx, "Readdir: %s -> %s", f.Name, fullPath)

			if f.IsDir {
				log.Infof(ctx, "Readdir: %s is a directory", f.Name)
				inode = r.NewInode(ctx, &FreighterRoot{Client: r.Client, IsDir: true, Repository: r.Repository, Target: r.Target, Path: fullPath}, fs.StableAttr{Mode: syscall.S_IFDIR, Ino: inodeId})
			} else {
				log.Infof(ctx, "Readdir: %s is a file", f.Name)
				inode = r.NewInode(ctx, &FreighterNode{Client: r.Client, Repository: r.Repository, Target: r.Target, Name: f.Name, Size: f.Size, Path: fullPath}, fs.StableAttr{Mode: 0755, Ino: inodeId})
			}
			r.Children[f.Name] = inode
			inodes = append(inodes, inode)
			filenames = append(filenames, f.Name)
			r.AddChild(f.Name, inode, false)
			log.Infof(ctx, "Readdir inode: %s -> %v (%s)", fullPath, *inode, inode.IsDir())
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
		//log.Infof(d.ctx, "Next inode: %s -> %v", fname, *inode)
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
	//log.Infof(d.ctx, "HasNext: %v", d.root.counter)
	if d.root.counter < len(d.filenames) {
		d.root.counter++
		return true
	}
	return false
}

func (r *FreighterRoot) PathTo(fname string) string {
	return fmt.Sprintf("%s/%s", r.Path, fname)
}

func (r *FreighterRoot) Mode() uint32 {
	if r.IsDir {
		return syscall.S_IFDIR | 0755
	} else {
		return 0755
	}
}

func (r *FreighterRoot) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	//log.Infof(ctx, "Lookup %v", r.PathTo(name))
	if i, ok := r.Children[name]; !ok {
		return nil, syscall.ENOENT
	} else {
		return i, 0
	}
}

func (r *FreighterRoot) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	//_, file := strings.Split(r.Path, "/")[0], strings.Split(r.Path, "/")[1]
	file := r.Path
	log.Infof(ctx, "OPEN %s:%s %s", r.Repository, r.Target, file)
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
	log.Infof(ctx, "OPEN %s:%s %s", f.Repository, f.Target, f.Path)
	if f.Data == nil {
		resp, err := f.Client.GetFile(ctx, &pb.FileRequest{Repository: f.Repository, Target: f.Target, Path: f.Path})
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
