package service

import (
	"context"
	"encoding/json"
	"log"
	"microservices/geo/internal/infrastructure/errors"
	"microservices/geo/internal/models"
	"microservices/geo/internal/modules/geo/storage"
	"strconv"

	"time"

	"github.com/ekomobile/dadata/v2"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/ptflp/gopubsub/queue"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var (
	defaultTimeout = time.Duration(1000 * time.Millisecond)
)

var (
	GeoControllerSearchAPIDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "geo_controller_search_api_request_duration_seconds",
		Help: "Request to API duration in seconds",
	})
)

func init() {
	prometheus.MustRegister(GeoControllerSearchAPIDuration)
}

type Geo struct {
	storage storage.GeoStorager
	logger  *zap.Logger
	limit   ratelimit.Limiter
	mq      queue.MessageQueuer
}

func NewGeo(storage storage.GeoStorager, logger *zap.Logger, limit ratelimit.Limiter, mq queue.MessageQueuer) Georer {
	return &Geo{storage: storage, logger: logger, limit: limit, mq: mq}
}

func (g *Geo) GeoCode(in GeoCodeIn) GeoCodeOut {
	return GeoCodeOut{
		Lat: in.Lat,
		Lng: in.Lng,
	}
}

func (g *Geo) SearchAddresses(ctx context.Context, in SearchAddressesIn) SearchAddressesOut {

	address, err := g.storage.Select(in.Query)
	if err != nil {

		res, err := g.searchFromAPI(ctx, in.Query)
		if err != nil {

			if err == errors.ErrRateLimitExceeded {
				//err = g.publishMessage(ctx)
				g.publishMessage(ctx)
			}
			log.Println("(g *Geo) SearchAddresses err", err)
			return SearchAddressesOut{
				Err: err,
			}
		}

		if err = g.storage.Insert(in.Query, res.Lat, res.Lon); err != nil {
			g.logger.Error("ошибка при добавлении данных в бд:", zap.Error(err))
		} else {
			g.logger.Info("Данные добавлены в бд")
		}

		return SearchAddressesOut{
			Address: res,
		}
	}

	return SearchAddressesOut{
		Address: address,
		Err:     nil,
	}
}

func (g *Geo) searchFromAPI(ctx context.Context, query string) (models.Address, error) {
	if !tryWithTimeout(g.limit, defaultTimeout) {
		return models.Address{}, errors.ErrRateLimitExceeded
	}

	startTime := time.Now()

	// api := dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
	// 	ApiKeyValue:    "d538755936a28def6bca48517dd287303cb0dae7",
	// 	SecretKeyValue: "81081aa1fa5ca90caa8a69b14947b5876f58b8db",
	// }))

	api := dadata.NewCleanApi()

	addresses, err := api.Address(ctx, query)
	if err != nil {
		return models.Address{}, err
	}

	durations := time.Since(startTime).Seconds()
	GeoControllerSearchAPIDuration.Observe(durations)

	res := models.Address{
		Lat: addresses[0].GeoLat,
		Lon: addresses[0].GeoLon,
	}

	return res, nil
}

func tryWithTimeout(limiter ratelimit.Limiter, timeout time.Duration) bool {
	done := make(chan struct{})

	go func() {
		limiter.Take()
		close(done)
	}()

	timer := time.NewTimer(timeout)

	select {
	case <-done:
		return true
	case <-timer.C:
		return false
	}
}

func (g *Geo) publishMessage(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)

	//_, claims, _ := jwtauth.FromContext(ctx)

	// email := claims["email"].(string)
	// phone := claims["phone"].(string)

	// g.logger.Info("EMAIL AND PHONE:", zap.String("Email", email), zap.String("Phone", phone))

	idRaw := md.Get("id")[0]

	id, err := strconv.Atoi(idRaw)
	if err != nil {
		g.logger.Error("convert id err", zap.Error(err))
	}

	msg := RateLimitMsg{
		ID: id,
	}

	msgJSON, _ := json.Marshal(msg)
	if err := g.mq.Publish("rate_limit", msgJSON); err != nil {
		g.logger.Error("publish message err", zap.Error(err))
		return err
	}
	return nil
}
