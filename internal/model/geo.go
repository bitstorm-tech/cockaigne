package model

import "fmt"

type Point struct {
	Lon float64
	Lat float64
}

func (p Point) ToWkt() string {
	return fmt.Sprintf("POINT(%f %f)", p.Lon, p.Lat)
}

type NominatimResponse struct {
	Lat     string `json:"lat"`
	Lon     string `json:"lon"`
	Display string `json:"display_name"`
}
