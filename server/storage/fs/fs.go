package fs

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnewart/freighter/server/storage/types"
	"zombiezen.com/go/log"
)

type DiskLayerFileStore struct {
	types.LayerStore
	root string
}

func NewDiskLayerFileStore(root string) (*DiskLayerFileStore, error) {
	return &DiskLayerFileStore{
		root: root,
	}, nil
}

func (s *DiskLayerFileStore) readFilesFromArchive(digest types.Digest) ([]types.FileRecord, error) {
	digestPath := s.getLayerPath(digest)
	log.Infof(context.Background(), "Reading files from layer: %s @ %s", digest, digestPath)
	files := make([]types.FileRecord, 0)
	if file, err := os.Open(digestPath); err == nil {
		archive, err := gzip.NewReader(file)

		if err != nil {
			return nil, err
		}

		tr := tar.NewReader(archive)

		if tr == nil {
			return nil, err
		}

		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			files = append(files, types.FileRecord{Name: hdr.Name, Size: hdr.Size, IsDir: false})
		}
	}

	return files, nil
}

func (s *DiskLayerFileStore) LayerFileNames(digest types.Digest) ([]string, error) {
	files, err := s.readFilesFromArchive(digest)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(files))
	for _, f := range files {
		names = append(names, f.Name)
	}
	return names, nil
}

func (s *DiskLayerFileStore) IngestLayer(ctx context.Context, digest types.Digest) ([]types.LayerFile, error) {
	if files, err := s.readFilesFromArchive(digest); err != nil {
		log.Errorf(context.Background(), "Error reading files from layer: %v", err)
		return nil, err
	} else {
		result := make([]types.LayerFile, 0, len(files))
		for _, r := range files {
			if !strings.HasPrefix(r.Name, "/") {
				r.Name = fmt.Sprintf("/%s", r.Name)
			}

			parts := strings.Split(r.Name, "/")
			dir := strings.Join(parts[:len(parts)-1], "/")

			if strings.HasSuffix(r.Name, "/") {
				r.IsDir = true
				if len(parts) > 2 {
					dir = strings.Join(parts[:len(parts)-2], "/")
				}
			}

			if dir == "" {
				dir = "/"
			}

			log.Infof(context.Background(), "Ingesting file: %s in '%s' (%v)", r.Name, dir, r.IsDir)

			result = append(result, types.LayerFile{
				LayerDigest: digest.String(),
				FilePath:    r.Name,
				Size:        r.Size,
				IsDir:       r.IsDir,
				Directory:   dir,
			})

		}

		return result, nil
	}
}

func (s *DiskLayerFileStore) readFile(digest types.Digest, filename string) ([]byte, error) {
	ctx := context.Background()
	log.Infof(ctx, "Fetching file %s from  %s", filename, digest)

	filename = strings.TrimPrefix(filename, "/")
	tarFile := s.getLayerPath(digest)
	file, err := os.Open(tarFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	archive, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer archive.Close()

	tr := tar.NewReader(archive)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.Name == filename {
			buf, err := ioutil.ReadAll(tr)
			if err != nil {
				return nil, err
			}
			return buf, nil
		}
	}
	return nil, fmt.Errorf("File not found: %s", filename)
}

func (s *DiskLayerFileStore) GetDirectoryTreeForLayer(digest types.Digest) ([]types.FileRecord, error) {
	return s.readFilesFromArchive(digest)
}

func (s *DiskLayerFileStore) DeleteLayer(digest types.Digest) error {
	return os.Remove(s.getLayerPath(digest))
}

func (s *DiskLayerFileStore) getLayerPath(digest types.Digest) string {
	return filepath.Join(s.root, digest.Algorithm, digest.Hash)
}

func (d *DiskLayerFileStore) GetLayer(digest types.Digest) (*types.Layer, error) {
	layerPath := d.getLayerPath(digest)
	log.Infof(context.Background(), "Layer not found in DB: %s - checking if exists on disk at: %s", digest, layerPath)
	if stat, err := os.Stat(layerPath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		layer := &types.Layer{
			Digest: digest.String(),
			Size:   stat.Size(),
		}
		return layer, nil
	}
}

func (s *DiskLayerFileStore) ReadFile(lf types.LayerFile) ([]byte, error) {
	return s.readFile(lf.Digest(), lf.FilePath)
}

func (s *DiskLayerFileStore) blobHashPath(algorithm string, h string) string {
	return filepath.Join(s.root, algorithm, h)
}

func (s *DiskLayerFileStore) StoreBlob(digest types.Digest, rc io.ReadCloser) error {
	f, err := os.CreateTemp(s.root, "upload-*")
	if err != nil {
		return err
	}

	if err := func() error {
		defer f.Close()
		_, err := io.Copy(f, rc)
		return err
	}(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(s.root, digest.Algorithm), os.ModePerm); err != nil {
		return err
	}

	layerPath := s.getLayerPath(digest)
	if _, err := os.Stat(layerPath); err == nil {
		log.Infof(context.TODO(), "Layer already exists: %s", digest)
		os.Remove(layerPath)
	}

	err = os.Rename(f.Name(), layerPath)
	if err != nil {
		log.Errorf(context.TODO(), "Error renaming file: %v", err)
		return err
	}

	return nil
}

func (s *DiskLayerFileStore) String() string {
	return fmt.Sprintf("DiskLayerFileStore: %s", s.root)
}

func (s *DiskLayerFileStore) GetLayerReader(digest types.Digest) (io.ReadCloser, error) {
	return os.Open(s.getLayerPath(digest))
}

var _ = (types.LayerStore)((*DiskLayerFileStore)(nil))
