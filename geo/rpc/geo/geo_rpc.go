package geo

import (
	"context"
	"microservices/geo/internal/modules/geo/service"

	pb "microservices/geo/geogrpc"
)

type GeoServiceRPC struct {
	geoService service.Georer
}

func NewGeoServiceRPC(geoService service.Georer) *GeoServiceRPC {
	return &GeoServiceRPC{geoService: geoService}
}

func (g *GeoServiceRPC) SearchAddresses(in service.SearchAddressesIn, out *service.SearchAddressesOut) error {
	ctx := context.Background()
	*out = g.geoService.SearchAddresses(ctx, in)
	return nil
}

func (g *GeoServiceRPC) GeoCode(in service.GeoCodeIn, out *service.GeoCodeOut) error {
	*out = g.geoService.GeoCode(in)

	return nil
}

type GeoServiceGRPC struct {
	geoService service.Georer
	pb.UnimplementedGeorerServer
}

func NewGeoServiceGRPC(geoService service.Georer) *GeoServiceGRPC {
	return &GeoServiceGRPC{
		geoService: geoService,
	}
}

func (g *GeoServiceGRPC) SearchAddresses(ctx context.Context, in *pb.SearchAddressesRequest) (*pb.SearchAddressesResponse, error) {
	address := g.geoService.SearchAddresses(ctx, service.SearchAddressesIn{Query: in.Query})

	response := pb.Address{
		Lat: address.Address.Lat,
		Lon: address.Address.Lon,
	}
	return &pb.SearchAddressesResponse{Address: &response}, address.Err
}

func (g *GeoServiceGRPC) GeoCode(ctx context.Context, in *pb.GeoCodeRequest) (*pb.GeoCodeResponse, error) {
	res := g.geoService.GeoCode(service.GeoCodeIn{Lat: in.Lat, Lng: in.Lng})

	return &pb.GeoCodeResponse{Lat: res.Lat, Lng: res.Lng}, res.Err
}
