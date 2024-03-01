package server

import (
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
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
			return nil, err
		}
	}

	return &layerFile, nil
}

func (d *DB) GetFilesForRepo(repo string, target string) ([]LayerFile, error) {
	var layers []Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Find(&layers).Error; err != nil {
		return nil, err
	}

	var layerFiles []LayerFile

	for _, l := range layers {
		digest := strings.Split(l.Digest, ":")[1]
		log.Infof(context.Background(), "Fetching files for layer: %s:%s %s", repo, target, digest)
		if err := d.db.Where("layer_digest = ?", digest).Find(&layerFiles).Error; err != nil {
			return nil, err
		}
	}

	return layerFiles, nil
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
