package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/gofiber/fiber/v2/log"
)

func GetPositionFromAddress(city string, zip int, street string, houseNumber string) (model.Point, error) {
	city = url.QueryEscape(city)
	street = url.QueryEscape(street)
	houseNumber = url.QueryEscape(houseNumber)
	query := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&street=%s,%s&city=%s&postalcode=%s", houseNumber, street, city, strconv.Itoa(zip))

	return pointFromQuery(query)
}

func GetPositionFromAddressFuzzy(address string) (model.Point, error) {
	return pointFromQuery("https://nominatim.openstreetmap.org/search?format=json&q=" + url.QueryEscape(address))
}

func pointFromQuery(query string) (model.Point, error) {
	log.Debugf("Querying nominatim: %s", query)
	res, err := http.Get(query)
	if err != nil {
		return model.Point{}, fmt.Errorf("can't get response from nominatim: %v", err)
	}

	defer res.Body.Close()
	var results []model.NominatimResponse
	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		return model.Point{}, fmt.Errorf("can't parse response from nominatim: %v", err)
	}

	if len(results) == 0 {
		return model.Point{}, fmt.Errorf("can't find location for address")
	}

	if len(results) > 1 {
		log.Warnf("Found more then one (%s) positions for address (use first one!)", len(results))
	}

	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return model.Point{}, fmt.Errorf("error parsing lon: %v", err)
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return model.Point{}, fmt.Errorf("error parsing lat: %v", err)
	}

	return model.Point{
		Lon: lon,
		Lat: lat,
	}, nil
}
