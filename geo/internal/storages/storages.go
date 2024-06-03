package storages

import (
	"microservices/geo/internal/db/adapter"
	"microservices/geo/internal/modules/geo/storage"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Storages struct {
	Geo storage.GeoStorager
}

func NewStorages(sqlAdapter *adapter.SQLAdapter, cache *redis.Client, logger *zap.Logger) *Storages {
	geoStorage := storage.NewGeoStorage(sqlAdapter)
	proxy := storage.NewGeoStorageProxy(geoStorage, cache, logger)

	return &Storages{
		Geo: proxy,
	}
}
