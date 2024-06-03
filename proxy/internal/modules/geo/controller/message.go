package controller

import "proxy/internal/models"

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type GeocodeResponse struct {
	Addresses []*models.Address `json:"addresses"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Addresses []*models.Address `json:"addresses"`
}
