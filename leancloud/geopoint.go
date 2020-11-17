package leancloud

import "math"

// GeoPoint represent
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// RadiansTo return the distance from this GeoPoint to another in radians
func (point *GeoPoint) RadiansTo(target *GeoPoint) float64 {
	d2r := math.Pi / 180.0
	lat1rad := point.Latitude * d2r
	long1rad := point.Longitude * d2r

	lat2rad := target.Latitude * d2r
	long2rad := target.Longitude * d2r

	deltaLat := lat1rad - lat2rad
	deltaLong := long1rad - long2rad

	latSinDelta := math.Sin(deltaLat / 2.0)
	longSinDelta := math.Sin(deltaLong / 2.0)

	a := (latSinDelta * longSinDelta) + (math.Cos(lat1rad) * math.Cos(lat2rad) * longSinDelta * longSinDelta)

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
