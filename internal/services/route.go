package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"task/internal/dto"
	"task/internal/entities"
	"task/internal/repositories"
)

type RouteService interface {
	Register(ctx context.Context, data dto.RegisterRouteRequestBody) error
	GetById(ctx context.Context, id int) (entities.Route, error)
	DeleteByIds(ctx context.Context, ids dto.DeleteRoutesRequestBody) error
}

type routeService struct {
	repo repositories.RouteRepo
}

func NewRouteService(repo repositories.RouteRepo) RouteService {
	return &routeService{repo: repo}
}

func (s *routeService) Register(ctx context.Context, data dto.RegisterRouteRequestBody) (err error) {
	route, err := dto.ToEntityModel(data)
	if err != nil {
		return fmt.Errorf("converting dto to entity model: %w", err)
	}

	err = s.repo.Register(ctx, route)
	if err != nil {
		return fmt.Errorf("route registration: %w", err)
	}

	return nil
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

// TODO: maybe implement in another way
func (s *routeService) DeleteByIds(ctx context.Context, ids dto.DeleteRoutesRequestBody) (err error) {
	errs := make([]error, 0)

	var wg sync.WaitGroup
	for _, id := range ids.RouteIDs {
		wg.Add(1)
		go func(routeID int) {
			inErr := s.repo.DeleteById(ctx, routeID)
			if inErr != nil {
				errs = append(errs, inErr)
			}
			wg.Done()
		}(id)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("deleting routes: %w", errors.Join(errs...))
	}

	return nil
}
