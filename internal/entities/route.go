package entities

type Route struct {
	RouteID   int
	RouteName string
	Load      float32
	CargoType string
	IsActual  bool
}
