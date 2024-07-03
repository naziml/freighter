package layers

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

	"github.com/johnewart/freighter/server/data"
	"zombiezen.com/go/log"
)

type DiskLayerFileStore struct {
	LayerFileStore
	RootPath string
	DB       *data.DB
}

type FileRecord struct {
	Name  string
	Size  int64
	IsDir bool
}

func NewDiskLayerFileStore(rootPath string, db *data.DB) *DiskLayerFileStore {
	return &DiskLayerFileStore{
		RootPath: rootPath,
		DB:       db,
	}
}

func (s *DiskLayerFileStore) getLayerPath(digest string) string {
	return filepath.Join(s.RootPath, "sha256", digest)
}

func (s *DiskLayerFileStore) readFilesFromArchive(digest string) ([]FileRecord, error) {
	digestPath := s.getLayerPath(digest)
	log.Infof(context.Background(), "Reading files from layer: %s @ %s", digest, digestPath)
	files := make([]FileRecord, 0)
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
			files = append(files, FileRecord{Name: hdr.Name, Size: hdr.Size, IsDir: false})
		}
	}

	return files, nil
}

func (s *DiskLayerFileStore) LayerFileNames(digest string) ([]string, error) {
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

func (s *DiskLayerFileStore) IngestFiles(digest string) error {
	if files, err := s.readFilesFromArchive(digest); err != nil {
		log.Errorf(context.Background(), "Error reading files from layer: %v", err)
		return err
	} else {
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

			log.Infof(context.Background(), "Ingesting file: %s in '%s' (%s)", r.Name, dir, r.IsDir)

			if err := s.DB.PutLayerFile(&data.LayerFile{
				LayerDigest: digest,
				FilePath:    r.Name,
				Size:        r.Size,
				IsDir:       r.IsDir,
				Directory:   dir,
			}); err != nil {
				log.Errorf(context.Background(), "Error creating layer file: %v", err)
				return err
			}
		}
		return nil
	}
}

func (s *DiskLayerFileStore) ReadFile(layerDigest string, filename string) ([]byte, error) {
	ctx := context.Background()
	log.Infof(ctx, "Fetching file %s from  %s", filename, layerDigest)

	filename = strings.TrimPrefix(filename, "/")
	tarFile := s.getLayerPath(layerDigest)
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

func (s *DiskLayerFileStore) GetDirectoryTree(layerDigest string) ([]FileRecord, error) {
	return s.readFilesFromArchive(layerDigest)
}

func (s *DiskLayerFileStore) DeleteLayer(hex string) error {

	if err := s.DB.DeleteLayer(hex); err != nil {
		log.Errorf(context.Background(), "Error deleting layer: %v", err)
		return err
	}

	return nil
}
