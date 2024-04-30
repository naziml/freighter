package layers

type LayerFileStore interface {
	StoreLayer(layerDigest string, data []byte) error
	LayerExists(repository, target, layerDigest string) bool
	ReadLayer(repository, target, layerDigest string) ([]byte, error)
}
