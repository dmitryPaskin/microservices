package controller

import (
	"encoding/json"
	"net/http"
	"proxy/internal/infrastructure/component"
	"proxy/internal/infrastructure/responder"
	"proxy/internal/models"
	"proxy/internal/modules/geo/service"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	geoControllerGeocodeRequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "geo_controller_geocode_requests_total",
		Help: "Total number of requests",
	})
	geoControllerSearchTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "geo_controller_search_requests_total",
		Help: "Total number of requests",
	})
	geoControllerGeocodeDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "geo_controller_geocode_duration_seconds",
		Help: "Request duration in seconds",
	})
	geoControllerSearchDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "geo_controller_search_duration_seconds",
		Help: "Request duration in seconds",
	})
)

func init() {
	prometheus.MustRegister(geoControllerGeocodeRequestsTotal)
	prometheus.MustRegister(geoControllerSearchTotal)
	prometheus.MustRegister(geoControllerGeocodeDuration)
	prometheus.MustRegister(geoControllerSearchDuration)
}

type Georer interface {
	Geocode(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	ApiHandler(w http.ResponseWriter, r *http.Request)
}

type GeoController struct {
	geo service.Georer
	responder.Responder
}

func NewGeoController(service service.Georer, components *component.Components) Georer {
	return &GeoController{
		geo:       service,
		Responder: components.Responder,
	}
}

func (g *GeoController) Geocode(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	geoControllerGeocodeRequestsTotal.Inc()

	var geocodeRequest GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&geocodeRequest); err != nil {
		g.ErrorBadRequest(w, err)
		return
	}

	geo := g.geo.GeoCode(service.GeoCodeIn{Lat: geocodeRequest.Lat, Lng: geocodeRequest.Lng})

	geocodeResponse := GeocodeResponse{
		Addresses: []*models.Address{{Lat: geo.Lat, Lon: geo.Lng}},
	}

	g.OutputJSON(w, geocodeResponse)

	duration := time.Since(startTime).Seconds()
	geoControllerGeocodeDuration.Observe(duration)
}

func (g *GeoController) Search(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	geoControllerSearchTotal.Inc()

	var searchRequest SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&searchRequest); err != nil {
		g.ErrorBadRequest(w, err)
		return
	}

	out := g.geo.SearchAddresses(r.Context(), service.SearchAddressesIn{Query: searchRequest.Query})
	if out.Err != nil {
		switch status.Code(out.Err) {
		case codes.ResourceExhausted:
			g.ErrorToManyRequests(w, out.Err)
			return
		default:
			g.ErrorInternal(w, out.Err)
			return
		}
	}

	searchResponse := SearchResponse{
		Addresses: []*models.Address{&out.Address},
	}

	g.OutputJSON(w, searchResponse)

	durations := time.Since(startTime).Seconds()
	geoControllerSearchDuration.Observe(durations)
}

func (g GeoController) ApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from API"))
}
