package services

import (
	"context"
	"fmt"
	"task/internal/dto"
	"task/internal/entities"
	"task/internal/repositories"
	"time"
)

type RouteService interface {
	Register(ctx context.Context, data dto.RegisterRouteRequestBody) (bool, error)
	GetById(ctx context.Context, id int) (entities.Route, error)
	DeleteByIds(ctx context.Context, ids dto.DeleteRoutesRequestBody) error
}

type routeService struct {
	repo repositories.RouteRepo
}

func NewRouteService(repo repositories.RouteRepo) RouteService {
	return &routeService{repo: repo}
}

func (s *routeService) Register(ctx context.Context, data dto.RegisterRouteRequestBody) (old bool, err error) {
	route, err := dto.ToEntityModel(data)
	if err != nil {
		return false, fmt.Errorf("converting dto to entity model: %w", err)
	}

	old, err = s.repo.Register(ctx, route)
	if err != nil {
		return false, fmt.Errorf("route registration: %w", err)
	}

	return old, nil
}

func (s *routeService) GetById(ctx context.Context, id int) (route entities.Route, err error) {
	if id < 0 {
		return entities.Route{}, fmt.Errorf("route id should be non-negative")
	}

	route, err = s.repo.GetById(ctx, id)
	if err != nil {
		return entities.Route{}, fmt.Errorf("getting route by id: %w", err)
	}

	return route, nil
}

func (s *routeService) DeleteByIds(ctx context.Context, ids dto.DeleteRoutesRequestBody) (err error) {
	for _, val := range ids.RouteIDs {
		if val < 0 {
			return fmt.Errorf("deleting routes: ids should be non-negative")
		}
	}

	go func() {
		// new context because after getting http response on this request original request context is cancelled
		delCtx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		err = s.repo.DeleteById(delCtx, ids.RouteIDs)
		if err != nil {
			// TODO: replace with log msg
			fmt.Printf("deleting routes: %v\n", err)
		}
	}()

	return nil
}
