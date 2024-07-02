package repositories

import "task/internal/entities"

type RouteRepo interface {
	Register(cargo entities.Route) error
	GetById(id int) (entities.Route, error)
	DeleteById(id int) error
	SetInactual(id int) error
}
