package model

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	Lon float64
	Lat float64
}

var PointCenterOfGermany = Point{
	Lon: 10.447683,
	Lat: 51.163361,
}

func NewPointFromString(pointString string) (Point, error) {
	lonLat := strings.Split(pointString, ",")

	lon, err := strconv.ParseFloat(lonLat[0], 64)
	if err != nil {
		return Point{}, err
	}

	lat, err := strconv.ParseFloat(lonLat[1], 64)
	if err != nil {
		return Point{}, err
	}

	return Point{
		Lon: lon,
		Lat: lat,
	}, nil
}

func (p Point) ToWkt() string {
	return fmt.Sprintf("POINT(%f %f)", p.Lon, p.Lat)
}

type NominatimResponse struct {
	Lat     string `json:"lat"`
	Lon     string `json:"lon"`
	Display string `json:"display_name"`
}
