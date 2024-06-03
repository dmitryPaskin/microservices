// Package classification Geoservice.
//
// Documentation of Geoservice API.
//
//		Version: 1.0.0
//		Schemes:
//		- http
//		BasePath: /
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Security:
//		- basic
//
//		SecurityDefinitions:
//		 Bearer:
//	   type: apiKey
//	   name: Authorization
//	   in: header
//
// swagger:meta
package docs

//go:generate swagger generate spec -o ../../public/swagger.json --scan-models
