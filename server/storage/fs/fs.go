package fs

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/base64"
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

func (s *DiskLayerFileStore) IngestLayer(ctx context.Context, digest types.Digest) ([]types.LayerFile, error) {
	digestPath := s.getLayerPath(digest)
	log.Infof(context.Background(), "Reading files from layer: %s @ %s", digest, digestPath)
	result := make([]types.LayerFile, 0)

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

			r := types.FileRecord{Name: hdr.Name, Size: hdr.Size, IsDir: false}

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

			lf := types.LayerFile{
				LayerDigest: digest.String(),
				FilePath:    r.Name,
				Size:        r.Size,
				IsDir:       r.IsDir,
				Directory:   dir,
			}

			log.Infof(ctx, "Ingesting file: %s", lf.FilePath)
			outfilePath := s.getPathForLayerFile(lf)
			log.Infof(ctx, "Writing file to: %s", outfilePath)

			outDir, _ := filepath.Split(outfilePath)
			if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
				log.Errorf(ctx, "Error creating directory: %v", err)
				return nil, err
			}

			if f, err := os.Create(outfilePath); err == nil {
				if err := func() error {
					defer f.Close()
					_, err := io.Copy(f, tr)
					return err
				}(); err != nil {
					log.Infof(ctx, "Error copying file data: %v", err)
				}
			} else {
				log.Infof(ctx, "Error creating file: %v", err)
			}

			result = append(result, lf)

		}

		return result, nil
	} else {
		log.Errorf(ctx, "Error opening file: %v", err)
		return nil, err
	}

}

func (s *DiskLayerFileStore) DeleteLayer(digest types.Digest) error {
	return os.Remove(s.getLayerPath(digest))
}

func (s *DiskLayerFileStore) getLayerPath(digest types.Digest) string {
	return filepath.Join(s.root, digest.Algorithm, digest.Hash)
}

func (s *DiskLayerFileStore) getPathForLayerFile(lf types.LayerFile) string {
	log.Infof(context.Background(), "Getting path for layer file: %s", lf.String())
	digest := lf.Digest()
	hasher := sha1.New()
	hasher.Write([]byte(lf.FilePath))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return filepath.Join(s.root, "extracted", digest.Algorithm, digest.Hash, sha)
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
	fname := s.getPathForLayerFile(lf)
	log.Infof(context.Background(), "Reading file: %s", fname)
	if f, err := os.Open(fname); err == nil {
		defer f.Close()
		return ioutil.ReadAll(f)
	} else {
		return nil, err
	}
}

func (s *DiskLayerFileStore) blobHashPath(algorithm string, h string) string {
	return filepath.Join(s.root, algorithm, h)
}

func (s *DiskLayerFileStore) StoreLayerBlob(digest types.Digest, rc io.ReadCloser) error {
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
