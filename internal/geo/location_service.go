package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
)

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

func GetPositionFromAddress(city string, zip int, street string, houseNumber string) (Point, error) {
	city = url.QueryEscape(city)
	street = url.QueryEscape(street)
	houseNumber = url.QueryEscape(houseNumber)
	query := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&street=%s,%s&city=%s&postalcode=%s", houseNumber, street, city, strconv.Itoa(zip))

	return pointFromQuery(query)
}

func GetPositionFromAddressFuzzy(address string) (Point, error) {
	return pointFromQuery("https://nominatim.openstreetmap.org/search?format=json&q=" + url.QueryEscape(address))
}

func pointFromQuery(query string) (Point, error) {
	log.Debugf("Querying nominatim: %s", query)
	res, err := http.Get(query)
	if err != nil {
		return Point{}, fmt.Errorf("can't get response from nominatim: %v", err)
	}

	defer res.Body.Close()
	var results []NominatimResponse
	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		return Point{}, fmt.Errorf("can't parse response from nominatim: %v", err)
	}

	if len(results) == 0 {
		return Point{}, fmt.Errorf("can't find location for address")
	}

	if len(results) > 1 {
		log.Warnf("Found more then one (%s) positions for address (use first one!)", len(results))
	}

	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return Point{}, fmt.Errorf("error parsing lon: %v", err)
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return Point{}, fmt.Errorf("error parsing lat: %v", err)
	}

	return Point{
		Lon: lon,
		Lat: lat,
	}, nil
}
