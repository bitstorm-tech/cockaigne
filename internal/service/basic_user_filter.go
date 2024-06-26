package service

import (
	"github.com/bitstorm-tech/cockaigne/internal/model"
)

type BasicUserFilter struct {
	Location             model.Point
	SearchRadiusInMeters int
	UseLocationService   bool
	SelectedCategories   []int
}

// TODO refactor this to a more robust solution
// This solution is very dangerous -> attackers could flood this map (login as basic user is not restricted yet) -> OOM
var basicUserFilters = map[string]*BasicUserFilter{}

func NewBasicUser(id string) {
	basicUserFilters[id] = &BasicUserFilter{
		Location:             model.PointCenterOfGermany,
		SearchRadiusInMeters: 500,
		UseLocationService:   false,
	}
}

func GetBasicUserFilter(id string) *BasicUserFilter {
	filter, ok := basicUserFilters[id]

	if ok {
		return filter
	}

	NewBasicUser(id)

	return &BasicUserFilter{
		Location:             model.PointCenterOfGermany,
		SearchRadiusInMeters: 500,
		UseLocationService:   false,
	}
}

func DeleteBasicUser(id string) {
	delete(basicUserFilters, id)
}

func GetBasicUserSpatialFilter(id string) SpatialDealFilter {
	basicUserFilter := GetBasicUserFilter(id)
	return RadiusDealFilter{
		Radius: basicUserFilter.SearchRadiusInMeters,
		Point:  basicUserFilter.Location,
	}
}
