package leancloud

import "math"

// GeoPoint contains location's latitude and longitude
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// RadiansTo return the distance from this GeoPoint to another in radians
func (point *GeoPoint) RadiansTo(target *GeoPoint) float64 {
	radius := math.Pi / 180.0
	startLatRadius := point.Latitude * radius
	startLongRadius := point.Longitude * radius

	endLatRadius := target.Latitude * radius
	endLongRadius := target.Longitude * radius

	deltaLat := startLatRadius - endLatRadius
	deltaLong := startLongRadius - endLongRadius

	latSinDelta := math.Sin(deltaLat / 2.0)
	longSinDelta := math.Sin(deltaLong / 2.0)

	a := (latSinDelta * longSinDelta) + (math.Cos(startLatRadius) * math.Cos(endLatRadius) * longSinDelta * longSinDelta)

	a = math.Min(1.0, a)

	return (2 * math.Asin(math.Sqrt(a)))
}

// KilometersTo return the distance from this GeoPoint to another in kilometers
func (point *GeoPoint) KilometersTo(target *GeoPoint) float64 {
	return (point.RadiansTo(target) * 6371.0)
}

// MilesTo return the distance from this GeoPoint to another in miles
func (point *GeoPoint) MilesTo(target *GeoPoint) float64 {
	return (point.RadiansTo(target) * 3958.8)
}
