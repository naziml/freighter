package data

import (
	"context"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"zombiezen.com/go/log"
)

type DB struct {
	db *gorm.DB
}

type LayerFile struct {
	LayerDigest string
	FilePath    string
	Size        int64
	IsDir       bool
	Directory   string
}

type Layer struct {
	ManifestID uint
	MediaType  string
	Size       int
	Digest     string
	Level      int
	Repository string
	Target     string
}

func (l *Layer) Hash() string {
	if strings.HasPrefix(l.Digest, "sha256:") {
		return strings.Split(l.Digest, ":")[1]
	}
	return l.Digest
}

type Manifest struct {
	ID              uint   `gorm:"primaryKey"`
	Repository      string `gorm:"index:idx_repotarget,unique"`
	Target          string `gorm:"index:idx_repotarget,unique"`
	SchemaVersion   int
	MediaType       string
	ConfigMediaType string
	ConfigDigest    string
	ConfigSize      int
	RawBlob         []byte
}

func NewDB(path string) *DB {
	db, _ := gorm.Open(sqlite.Open(path), &gorm.Config{})
	d := &DB{
		db: db,
	}

	d.Migrate()
	return d
}

func (d *DB) Migrate() {
	d.db.AutoMigrate(&Manifest{})
	d.db.AutoMigrate(&Layer{})
	d.db.AutoMigrate(&LayerFile{})
}

func (d *DB) GetManifest(repo string, target string) (*Manifest, error) {
	var manifest Manifest

	if err := d.db.Where("repository = ? AND target = ?", repo, target).First(&manifest).Error; err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (d *DB) PutManifest(m Manifest) (Manifest, error) {

	if err := d.db.Create(&m).Error; err != nil {
		return m, err
	}

	return m, nil
}

func (d *DB) ManifestsForRepo(repo string) ([]Manifest, error) {
	var manifests []Manifest

	if err := d.db.Where("repository = ?", repo).Find(&manifests).Error; err != nil {
		return nil, err
	}

	return manifests, nil
}

func (d *DB) ManifestExists(repo string, target string) bool {
	var manifest Manifest

	if err := d.db.Where("repository = ? AND target = ?", repo, target).First(&manifest).Error; err != nil {
		return false
	}

	return true
}

func (d *DB) ListRepositories() []string {
	var repos []string

	if err := d.db.Model(&Manifest{}).Select("repository").Group("repository").Find(&repos).Error; err != nil {
		return nil
	}

	return repos
}

func (d *DB) DeleteManifest(repo string, target string) error {
	if err := d.db.Where("repository = ? AND target = ?", repo, target).Delete(&Manifest{}).Error; err != nil {
		return err
	}

	return nil
}

func (d *DB) GetFileLayer(repo string, target string, filePath string) (*LayerFile, error) {

	var layers []Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Order("level desc").Find(&layers).Error; err != nil {
		return nil, err
	}

	var layerFile LayerFile

	for _, l := range layers {
		if err := d.db.Where("layer_digest = ? AND file_path = ?", l.Hash(), filePath).First(&layerFile).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return nil, err
		}
	}

	return &layerFile, nil
}

func (d *DB) GetDirectoryTreeForRepo(repo string, target string) []string {
	var layers []Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Find(&layers).Error; err != nil {
		log.Errorf(context.Background(), "Error getting layers: %v", err)
		return nil
	}

	var layerDigests []string
	for _, l := range layers {
		log.Infof(context.Background(), "Fetching directory tree for layer: %s:%s %s", repo, target, l.Digest)
		layerDigests = append(layerDigests, l.Hash())
	}

	var layerfiles []LayerFile

	if err := d.db.Where("layer_digest IN ?", layerDigests).Distinct("directory").Find(&layerfiles).Error; err != nil {
		log.Errorf(context.Background(), "Error getting directory tree: %v", err)
		return nil
	}

	log.Infof(context.Background(), "Found %d directories", len(layerfiles))

	var directories []string
	for _, f := range layerfiles {
		directories = append(directories, f.Directory)
	}
	return directories
}

func (d *DB) GetFilesForRepo(repo string, target string, path string) ([]LayerFile, error) {
	var layers []Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Find(&layers).Error; err != nil {
		return nil, err
	}

	log.Infof(context.Background(), "Fetching files for %s:%s with %d layers", repo, target, len(layers))

	var layerFiles []LayerFile

	filemap := make(map[string]LayerFile)

	for _, l := range layers {
		digest := l.Hash()
		log.Infof(context.Background(), "Fetching files for layer: %s:%s %s", repo, target, digest)
		if path != "" {
			log.Infof(context.Background(), "Fetching files for layer: %s:%s %s in %s", repo, target, digest, path)
			if err := d.db.Where("layer_digest = ? AND directory = ?", digest, path).Find(&layerFiles).Error; err != nil {
				return nil, err
			}
		} else {
			if err := d.db.Where("layer_digest = ?", digest).Find(&layerFiles).Error; err != nil {
				return nil, err
			}
		}

		for _, f := range layerFiles {
			filemap[f.FilePath] = f
		}
	}

	var files []LayerFile
	for _, f := range filemap {
		files = append(files, f)
	}
	return files, nil
}

func (d *DB) GetLayerFiles(digest string) ([]LayerFile, error) {
	var layerFiles []LayerFile

	if err := d.db.Where("layer_digest = ?", digest).Find(&layerFiles).Error; err != nil {
		return nil, err
	}

	return layerFiles, nil
}

func (d *DB) PutLayerFile(lf *LayerFile) error {
	if err := d.db.Create(lf).Error; err != nil {
		return err
	}

	return nil
}

func (d *DB) PutLayer(l *Layer) error {
	if err := d.db.Create(l).Error; err != nil {
		return err
	}

	return nil
}

func (d *DB) DeleteLayer(digest string) error {
	if err := d.db.Where("digest = ?", digest).Delete(&Layer{}).Error; err != nil {
		return err
	}

	if err := d.db.Where("layer_digest = ?", digest).Delete(&LayerFile{}).Error; err != nil {
		return err
	}

	return nil
}
