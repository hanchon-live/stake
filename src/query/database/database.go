package database

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type Database struct {
	Cache *ttlcache.Cache[string, string]
}

func NewDatabase() *Database {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, string](30*time.Second),
		ttlcache.WithDisableTouchOnHit[string, string](),
	)
	go cache.Start()

	return &Database{
		Cache: cache,
	}
}
