package storage

import (
	"context"
	"encoding/json"
	"microservices/geo/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type GeoStorageProxy struct {
	storage *GeoStorage
	cache   *redis.Client
	logger  *zap.Logger
}

func NewGeoStorageProxy(storage *GeoStorage, cache *redis.Client, logger *zap.Logger) GeoStorager {
	return &GeoStorageProxy{
		storage: storage,
		cache:   cache,
		logger:  logger,
	}
}

func (g *GeoStorageProxy) Select(query string) (models.Address, error) {
	startTime := time.Now()

	data, err := g.cache.Get(context.Background(), query).Result()

	duration := time.Since(startTime).Seconds()
	GeoControllerSearchCacheDuration.Observe(duration)

	if err == nil {
		var address models.Address

		err = json.Unmarshal([]byte(data), &address)
		if err == nil {
			g.logger.Info("данные получены из кэша")
			return address, nil
		}
		g.logger.Error("ошибка при разборе данных из кэша: ", zap.Error(err))
	}

	if err == redis.Nil {
		g.logger.Info("данных нет в кэшэ")
	} else if err != nil {
		g.logger.Error("ошибка при полученни данных из кэша: ", zap.Error(err))
	}

	address, err := g.storage.Select(query)
	if err != nil {
		return models.Address{}, err
	}

	g.logger.Info("данные получены из бд")

	err = g.cache.Set(context.Background(), query, address, 5*time.Minute).Err()
	if err != nil {
		g.logger.Error("ошибка при сохранении данных в кэш", zap.Error(err))
	} else {
		g.logger.Info("данные записаны в кэш")
	}

	return address, nil
}

func (g *GeoStorageProxy) Insert(query string, lat string, lon string) error {
	address := models.Address{
		Lat: lat,
		Lon: lon,
	}

	err := g.cache.Set(context.Background(), query, address, 5*time.Minute).Err()
	if err != nil {
		g.logger.Error("ошибка при записи данных кэш", zap.Error(err))
	} else {
		g.logger.Info("данные записаны в кэш")
	}

	return g.storage.Insert(query, lat, lon)
}
