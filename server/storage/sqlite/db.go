package sqlite

import (
	"context"
	"os"
	"time"

	golog "log"

	"github.com/johnewart/freighter/server/storage/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"zombiezen.com/go/log"
)

type DBMetadataStore struct {
	types.MetadataStore
	db *gorm.DB
}

func NewDBMetadataStore(path string) (*DBMetadataStore, error) {
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

	d := &DBMetadataStore{
		db: db,
	}

	d.Migrate()
	return d, nil
}

func (d *DBMetadataStore) Migrate() {
	d.db.AutoMigrate(&types.Manifest{})
	d.db.AutoMigrate(&types.Layer{})
	d.db.AutoMigrate(&types.LayerFile{})
}

func (d *DBMetadataStore) Close() {
	sqlDB, _ := d.db.DB()
	sqlDB.Close()
}

func (d *DBMetadataStore) GetLayer(digest types.Digest) (*types.Layer, error) {
	var layer types.Layer

	if err := d.db.Where("digest = ?", digest.String()).First(&layer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &layer, nil
}

func (d *DBMetadataStore) GetManifest(repo string, target string) (*types.Manifest, error) {
	var manifest types.Manifest

	if err := d.db.Where("repository = ? AND target = ?", repo, target).First(&manifest).Error; err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (d *DBMetadataStore) PutManifest(m types.Manifest) (types.Manifest, error) {

	if err := d.db.Create(&m).Error; err != nil {
		return m, err
	}

	return m, nil
}

func (d *DBMetadataStore) ManifestsForRepo(repo string) ([]types.Manifest, error) {
	var manifests []types.Manifest

	if err := d.db.Where("repository = ?", repo).Find(&manifests).Error; err != nil {
		return nil, err
	}

	return manifests, nil
}

func (d *DBMetadataStore) ManifestExists(repo string, target string) bool {
	var manifest types.Manifest

	if err := d.db.Where("repository = ? AND target = ?", repo, target).First(&manifest).Error; err != nil {
		return false
	}

	return true
}

func (d *DBMetadataStore) ListRepositories() []string {
	var repos []string

	if err := d.db.Model(&types.Manifest{}).Select("repository").Group("repository").Find(&repos).Error; err != nil {
		return nil
	}

	return repos
}

func (d *DBMetadataStore) DeleteManifest(repo string, target string) error {
	if err := d.db.Where("repository = ? AND target = ?", repo, target).Delete(&types.Manifest{}).Error; err != nil {
		return err
	}

	return nil
}

func (d *DBMetadataStore) GetFileLayer(repo string, target string, filePath string) (*types.LayerFile, error) {

	var layers []types.Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Order("level desc").Find(&layers).Error; err != nil {
		return nil, err
	}

	var layerFile types.LayerFile

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

func (d *DBMetadataStore) GetDirectoryTreeForRepo(repo string, target string) []string {
	var layers []types.Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Find(&layers).Error; err != nil {
		log.Errorf(context.Background(), "Error getting layers: %v", err)
		return nil
	}

	var layerDigests []string
	for _, l := range layers {
		log.Infof(context.Background(), "Fetching directory tree for layer: %s:%s %s", repo, target, l.Digest)
		layerDigests = append(layerDigests, l.Digest)
	}

	var layerfiles []types.LayerFile

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

func (d *DBMetadataStore) GetFilesForRepo(repo string, target string, path string) ([]types.LayerFile, error) {
	var layers []types.Layer

	if err := d.db.Where("repository = ? AND target = ?", repo, target).Find(&layers).Error; err != nil {
		return nil, err
	}

	log.Infof(context.Background(), "Fetching files for %s:%s with %d layers", repo, target, len(layers))

	var layerFiles []types.LayerFile

	filemap := make(map[string]types.LayerFile)

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

	var files []types.LayerFile
	for _, f := range filemap {
		files = append(files, f)
	}
	return files, nil
}

func (d *DBMetadataStore) GetLayerFiles(digest types.Digest) ([]types.LayerFile, error) {
	var layerFiles []types.LayerFile

	if err := d.db.Where("layer_digest = ?", digest.String()).Find(&layerFiles).Error; err != nil {
		return nil, err
	}

	return layerFiles, nil
}

func (d *DBMetadataStore) PutLayerFile(lf *types.LayerFile) error {
	if err := d.db.Create(lf).Error; err != nil {
		return err
	}

	return nil
}

func (d *DBMetadataStore) PutLayer(l *types.Layer) error {
	if err := d.db.Create(l).Error; err != nil {
		return err
	}

	return nil
}

func (d *DBMetadataStore) DeleteLayer(digest types.Digest) error {
	if err := d.db.Where("digest = ?", digest).Delete(&types.Layer{}).Error; err != nil {
		return err
	}

	if err := d.db.Where("layer_digest = ?", digest).Delete(&types.LayerFile{}).Error; err != nil {
		return err
	}

	return nil
}
