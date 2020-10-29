package leancloud

import (
	"fmt"
	"reflect"
)

type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func decodeGeoPoint(v map[string]interface{}) (*GeoPoint, error) {
	if v["__type"] != "GeoPint" {
		return nil, fmt.Errorf("want type GeoPoint but %s", v["__type"])
	}
	latitude, ok := v["latitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("latitude want type float64 but %v", reflect.TypeOf(v["latitude"]))
	}
	longitude, ok := v["longitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("longitude want type float64 but %v", reflect.TypeOf(v["longitude"]))
	}
	return &GeoPoint{
		Latitude: latitude,
		Longitude: longitude,
	}, nil
}