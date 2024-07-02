package dto

import (
	"fmt"
	"task/internal/entities"
)

const eps = 1e-6

type RegisterRouteRequestBody struct {
	RouteID   int     `json:"route_id"`
	RouteName string  `json:"route_name"`
	Load      float32 `json:"load"`
	CargoType string  `json:"cargo_type"`
}

type DeleteRoutesRequestBody struct {
	RouteIDs []int `json:"route_ids"`
}

func ToEntityModel(data RegisterRouteRequestBody) (route entities.Route, err error) {
	if data.RouteID < 0 {
		return entities.Route{}, fmt.Errorf("route id should be non-negative")
	}

	if data.RouteName == "" {
		return entities.Route{}, fmt.Errorf("route name should not be empty")
	}

	if data.Load < 0.0 {
		return entities.Route{}, fmt.Errorf("load should be non-negative")
	}

	if data.CargoType == "" {
		return entities.Route{}, fmt.Errorf("cargo type should not be empty")
	}

	return entities.Route{
		RouteID:   data.RouteID,
		RouteName: data.RouteName,
		Load:      data.Load,
		CargoType: data.CargoType,
	}, nil
}
