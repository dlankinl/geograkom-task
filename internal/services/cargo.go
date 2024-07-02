package services

import (
	"errors"
	"fmt"
	"sync"
	"task/internal/dto"
	"task/internal/entities"
	"task/internal/repositories"
)

type RouteService interface {
	Register(data dto.RegisterRouteRequestBody) error
	GetById(id int) (entities.Route, error)
	DeleteByIds(ids dto.DeleteRoutesRequestBody) error
}

type routeService struct {
	repo repositories.RouteRepo
}

func NewRouteService(repo repositories.RouteRepo) RouteService {
	return &routeService{repo: repo}
}

func (s *routeService) Register(data dto.RegisterRouteRequestBody) (err error) {
	route, err := dto.ToEntityModel(data)
	if err != nil {
		return fmt.Errorf("converting dto to entity model: %w", err)
	}

	err = s.repo.Register(route)
	if err != nil {
		return fmt.Errorf("route registration: %w", err)
	}

	return nil
}

func (s *routeService) GetById(id int) (route entities.Route, err error) {
	route, err = s.repo.GetById(id)
	if err != nil {
		return entities.Route{}, fmt.Errorf("getting route by id: %w", err)
	}

	return route, nil
}

func (s *routeService) DeleteByIds(ids dto.DeleteRoutesRequestBody) (err error) {
	errs := make([]error, 0)

	var wg sync.WaitGroup
	for _, id := range ids.RouteIDs {
		wg.Add(1)
		go func(routeID int) {
			inErr := s.repo.DeleteById(routeID)
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
