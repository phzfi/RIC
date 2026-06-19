package cache

import (
	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/logging"
)

func NewS3Cache(cfg config.CacheConfig, policy Policy) (*Cache, error) {
	store, err := NewS3Store(cfg)
	if err != nil {
		return nil, err
	}

	logging.Debugf("S3 cache create: bucket=%s, prefix=%s", cfg.S3Bucket, cfg.S3Prefix)

	return &Cache{
		maxMemory: 0,
		policy:    policy,
		storer:    store,
	}, nil
}