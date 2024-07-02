package repositories

import (
	"context"
	"task/internal/entities"
)

//go:generate mockgen -source=route.go -destination=../mocks/route.go -package=mocks
type RouteRepo interface {
	Register(ctx context.Context, cargo entities.Route) error
	GetById(ctx context.Context, id int) (entities.Route, error)
	DeleteById(ctx context.Context, id int) error
	SetInactual(ctx context.Context, id int) error
}
