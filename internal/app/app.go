package app

import (
	"github.com/jackc/pgx/v5"
	"task/internal/repositories"
	"task/internal/services"
)

type App struct {
	Svc services.RouteService
}

func NewApp(db *pgx.Conn) *App {
	repo := repositories.NewRouteRepo(db)
	svc := services.NewRouteService(repo)

	return &App{Svc: svc}
}
