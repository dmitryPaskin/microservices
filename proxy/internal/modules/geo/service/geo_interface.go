package service

import (
	"context"
	"proxy/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=Georer
type Georer interface {
	SearchAddresses(ctx context.Context, in SearchAddressesIn) SearchAddressesOut
	GeoCode(in GeoCodeIn) GeoCodeOut
}

type GeoCodeIn struct {
	Lat string
	Lng string
}

type GeoCodeOut struct {
	Lat string
	Lng string
	Err error
}

type SearchAddressesIn struct {
	Query string
}

type SearchAddressesOut struct {
	Address models.Address
	Err     error
}
