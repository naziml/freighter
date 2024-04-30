package layers

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type CachedLayerFileStore struct {
	LayerFileStore
	storage *DiskLayerFileStore
	redis   *redis.Client
}

func NewCachedLayerFileStore(storage *DiskLayerFileStore, redis *redis.Client) *CachedLayerFileStore {
	return &CachedLayerFileStore{
		storage: storage,
		redis:   redis,
	}
}

func (c *CachedLayerFileStore) ReadFile(ctx context.Context, repo string, target string, path string) ([]byte, error) {
	key := c.redisKey(repo, target, path)
	if data, err := c.redis.Get(ctx, key).Bytes(); err == nil {
		return data, nil
	}
	data, err := c.storage.ReadFile(repo, target, path)
	if err != nil {
		return nil, err
	}
	c.redis.Set(ctx, key, data, 0)
	return data, nil
}

func (c *CachedLayerFileStore) GetDirectoryTree(ctx context.Context, repo string, target string) []string {
	return c.storage.GetDirectoryTree(repo, target)
}

func (c *CachedLayerFileStore) ListFiles(ctx context.Context, repo string, target string, path string) ([]FileRecord, error) {
	return c.storage.ListFiles(repo, target, path)
}

func (c *CachedLayerFileStore) redisKey(repo string, target string, path string) string {
	return repo + ":" + target + ":" + path
}
