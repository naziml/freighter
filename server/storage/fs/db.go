package fs

import (
	"context"
	"os"
	"time"

	golog "log"

	"github.com/johnewart/freighter/server/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"zombiezen.com/go/log"
)

type DB struct {
	db *gorm.DB
}

func NewDB(path string) (*DB, error) {
	newLogger := logger.New(
		golog.New(os.Stdout, "\r\n", golog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,  // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,  // Don't include params in the SQL log
			Colorful:                  false, // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, err
	}

	d := &DB{
		db: db,
	}

	d.Migrate()
	return d, nil
}

func (d *DB) Migrate() {
	d.db.AutoMigrate(&storage.Manifest{})
	d.db.AutoMigrate(&storage.Layer{})
	d.db.AutoMigrate(&storage.LayerFile{})
}

func (d *DB) Close() {
	sqlDB, _ := d.db.DB()
	sqlDB.Close()
}

func (d *DB) GetLayer(digest storage.Digest) (*storage.Layer, error) {
	var layer storage.Layer

	if err := d.db.Where("digest = ?", digest.String()).First(&layer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &layer, nil
}

func (d *DB) GetManifest(repo string, target string) (*storage.Manifest, error) {
	var manifest storage.Manifest

	if err := d.db.Where("repository = ? AND target = ?", repo, target).First(&manifest).Error; err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (d *DB) PutManifest(m storage.Manifest) (storage.Manifest, error) {

	if err := d.db.Create(&m).Error; err != nil {
		return m, err
	}

	return m, nil
}

func (d *DB) ManifestsForRepo(repo string) ([]storage.Manifest, error) {
	var manifests []storage.Manifest

	if err := d.db.Where("repository = ?", repo).Find(&manifests).Error; err != nil {
		return nil, err
	}

	return manifests, nil
}

func (d *DB) ManifestExists(repo string, target string) bool {
	var manifest storage.Manifest

	if err := d.db.Where("repository = ? AND target = ?", repo, target).First(&manifest).Error; err != nil {
		return false
	}

	return true
}

func (d *DB) ListRepositories() []string {
	var repos []string

	if err := d.db.Model(&storage.Manifest{}).Select("repository").Group("repository").Find(&repos).Error; err != nil {
		return nil
	}

	return repos
}

func (d *DB) DeleteManifest(repo string, target string) error {
	if err := d.db.Where("repository = ? AND target = ?", repo, target).Delete(&storage.Manifest{}).Error; err != nil {
		return err
	}

	return nil
}

func (d *DB) GetFileLayer(repo string, target string, filePath string) (*storage.LayerFile, error) {

	var layers []storage.Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Order("level desc").Find(&layers).Error; err != nil {
		return nil, err
	}

	var layerFile storage.LayerFile

	for _, l := range layers {
		if err := d.db.Where("layer_digest = ? AND file_path = ?", l.GetDigest(), filePath).First(&layerFile).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return nil, err
		}
	}

	return &layerFile, nil
}

func (d *DB) GetDirectoryTreeForRepo(repo string, target string) []string {
	var layers []storage.Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Find(&layers).Error; err != nil {
		log.Errorf(context.Background(), "Error getting layers: %v", err)
		return nil
	}

	var layerDigests []string
	for _, l := range layers {
		log.Infof(context.Background(), "Fetching directory tree for layer: %s:%s %s", repo, target, l.Digest)
		layerDigests = append(layerDigests, l.Digest)
	}

	var layerfiles []storage.LayerFile

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

func (d *DB) GetFilesForRepo(repo string, target string, path string) ([]storage.LayerFile, error) {
	var layers []storage.Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Find(&layers).Error; err != nil {
		return nil, err
	}

	log.Infof(context.Background(), "Fetching files for %s:%s with %d layers", repo, target, len(layers))

	var layerFiles []storage.LayerFile

	filemap := make(map[string]storage.LayerFile)

	for _, l := range layers {
		log.Infof(context.Background(), "Fetching files for layer: %s:%s %s", repo, target, l.Digest)
		if path != "" {
			log.Infof(context.Background(), "Fetching files for layer: %s:%s %s in %s", repo, target, l.Digest, path)
			if err := d.db.Where("layer_digest = ? AND directory = ?", l.Digest, path).Find(&layerFiles).Error; err != nil {
				return nil, err
			}
		} else {
			if err := d.db.Where("layer_digest = ?", l.Digest).Find(&layerFiles).Error; err != nil {
				return nil, err
			}
		}

		for _, f := range layerFiles {
			filemap[f.FilePath] = f
		}
	}

	var files []storage.LayerFile
	for _, f := range filemap {
		files = append(files, f)
	}
	return files, nil
}

func (d *DB) GetLayerFiles(digest storage.Digest) ([]storage.LayerFile, error) {
	var layerFiles []storage.LayerFile

	if err := d.db.Where("layer_digest = ?", digest.String()).Find(&layerFiles).Error; err != nil {
		return nil, err
	}

	return layerFiles, nil
}

func (d *DB) PutLayerFile(lf *storage.LayerFile) error {
	if err := d.db.Create(lf).Error; err != nil {
		return err
	}

	return nil
}

func (d *DB) PutLayer(l *storage.Layer) error {
	if err := d.db.Create(l).Error; err != nil {
		return err
	}

	return nil
}

func (d *DB) DeleteLayer(digest storage.Digest) error {
	if err := d.db.Where("digest = ?", digest).Delete(&storage.Layer{}).Error; err != nil {
		return err
	}

	if err := d.db.Where("layer_digest = ?", digest).Delete(&storage.LayerFile{}).Error; err != nil {
		return err
	}

	return nil
}
