package storage

import (
	"microservices/geo/internal/db/adapter"
	"microservices/geo/internal/models"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	GeoControllerSearchCacheDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "geo_controller_search_cache_request_duration_seconds",
		Help: "Request to cache duration in seconds",
	})
	GeoControllerSearchDBDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "geo_controller_search_db_request_duration_seconds",
		Help: "Request to db duration in seconds",
	})
)

func init() {
	prometheus.MustRegister(GeoControllerSearchCacheDuration)
	prometheus.MustRegister(GeoControllerSearchDBDuration)
}

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=GeoStorager
type GeoStorager interface {
	Select(query string) (models.Address, error)
	Insert(query, lat, lon string) error
}

type GeoStorage struct {
	adapter adapter.SQLAdapterer
}

func NewGeoStorage(sqlAdapter adapter.SQLAdapterer) *GeoStorage {
	return &GeoStorage{
		adapter: sqlAdapter,
	}
}

func (g *GeoStorage) Select(query string) (models.Address, error) {
	startTime := time.Now()

	address, err := g.adapter.Select(query)

	duration := time.Since(startTime).Seconds()
	GeoControllerSearchDBDuration.Observe(duration)

	return address, err
}

func (g *GeoStorage) Insert(query, lat, lon string) error {
	return g.adapter.Insert(query, lat, lon)
}
