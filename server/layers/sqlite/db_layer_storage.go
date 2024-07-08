package sqlite

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

	"github.com/johnewart/freighter/server/layers"
	"zombiezen.com/go/log"
)

type DiskLayerFileStore struct {
	layers.LayerFileStore
	RootPath string
	DB       *DB
}

type FileRecord struct {
	Name  string
	Size  int64
	IsDir bool
}

func NewDiskLayerFileStore(rootPath string, db *DB) *DiskLayerFileStore {
	return &DiskLayerFileStore{
		RootPath: rootPath,
		DB:       db,
	}
}

func (s *DiskLayerFileStore) getLayerPath(digest string) string {
	return filepath.Join(s.RootPath, "sha256", digest)
}

func (s *DiskLayerFileStore) IngestFiles(digest string) error {
	digestPath := s.getLayerPath(digest)
	log.Infof(context.Background(), "Reading files from layer: %s @ %s", digest, digestPath)
	if file, err := os.Open(digestPath); err != nil {
		return err
	} else {
		archive, err := gzip.NewReader(file)

		if err != nil {
			return err
		}

		tr := tar.NewReader(archive)

		if tr == nil {
			return err
		}

		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			r := FileRecord{Name: hdr.Name, Size: hdr.Size, IsDir: false}

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

			if !r.IsDir {

				buf := make([]byte, hdr.Size)
				_, err = tr.Read(buf)
				if err != nil {
					log.Errorf(context.Background(), "Error reading file: %v", err)
					continue
				}

				if err := s.DB.PutLayerFile(&LayerFile{
					LayerDigest: digest,
					LayerIndex:  -1,
					Repository:  "",
					Target:      "",
					FilePath:    r.Name,
					Size:        r.Size,
					IsDir:       r.IsDir,
					Directory:   dir,
					Data:        buf,
				}); err != nil {
					log.Errorf(context.Background(), "Error creating layer file: %v", err)
					return err
				}
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

	dirNames := s.DB.GetDirectoryTreeForLayer(layerDigest)
	fileRecords := make([]FileRecord, 0)
	for _, d := range dirNames {
		fileRecords = append(fileRecords, FileRecord{Name: d, Size: 0, IsDir: true})
		log.Infof(context.Background(), "Directory: %s", d)
	}
	return fileRecords, nil
}

func (s *DiskLayerFileStore) DeleteLayer(hex string) error {

	if err := s.DB.DeleteLayer(hex); err != nil {
		log.Errorf(context.Background(), "Error deleting layer: %v", err)
		return err
	}

	return nil
}
