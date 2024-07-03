package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"task/internal/entities"
)

//go:generate mockgen -source=route.go -destination=../mocks/route.go -package=mocks
type RouteRepo interface {
	Register(ctx context.Context, route entities.Route) (bool, error)
	GetById(ctx context.Context, id int) (entities.Route, error)
	DeleteById(ctx context.Context, ids []int) error
}

type routeRepo struct {
	db *pgx.Conn
}

func NewRouteRepo(db *pgx.Conn) RouteRepo {
	return &routeRepo{
		db: db,
	}
}

func (r *routeRepo) Register(ctx context.Context, route entities.Route) (old bool, err error) {
	err = r.db.QueryRow(
		ctx,
		`with try as (
				insert into routes(route_id, route_name, load, cargo_type)
					values($1, $2, $3, $4)
					on conflict(route_id) do update set
						is_actual = false
					returning (xmax = 0) as inserted
			), new_id as (
				select max(route_id) + 1 as route_id
				from routes
			), insert_new as (
				insert into routes(route_id, route_name, load, cargo_type)
					select new_id.route_id, $2, $3, $4
					from new_id
					where not exists (select 1 from try where inserted)
					returning true as inserted_new_id
			)
			select
				case
					when exists (select 1 from insert_new) then true
						else false
				end as result`,
		route.RouteID,
		route.RouteName,
		route.Load,
		route.CargoType,
	).Scan(&old)
	if err != nil {
		return false, fmt.Errorf("register route: %w", err)
	}

	return old, nil
}

func (r *routeRepo) GetById(ctx context.Context, id int) (route entities.Route, err error) {
	err = r.db.QueryRow(
		ctx,
		`select
    			route_id, 
    			route_name, 
    			load, 
       			cargo_type,
       			is_actual
			from routes
			where route_id=$1`,
		id,
	).Scan(
		&route.RouteID,
		&route.RouteName,
		&route.Load,
		&route.CargoType,
		&route.IsActual,
	)
	if err != nil {
		return entities.Route{}, fmt.Errorf("getting route by id: %w", err)
	}

	return route, nil
}

func (r *routeRepo) DeleteById(ctx context.Context, ids []int) (err error) {
	_, err = r.db.Exec(
		ctx,
		`delete from routes where route_id = any($1)`,
		ids,
	)
	if err != nil {
		return fmt.Errorf("deleting route by id: %w", err)
	}

	return nil
}
