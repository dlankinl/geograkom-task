package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"task/internal/entities"
	"task/internal/integration_tests"
	"testing"
)

// for running tests use "go test -cover ./..." from root

const eps = 1e-6

var testDbInstance *pgx.Conn

func TestMain(m *testing.M) {
	testDB := integration_tests.SetupTestDatabase()
	defer testDB.TearDown()
	testDbInstance = testDB.DbInstance
	err := integration_tests.SeedTestData(testDbInstance)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestDeleteById(t *testing.T) {
	repo := NewRouteRepo(testDbInstance)

	testCases := []struct {
		name    string
		ids     []int
		wantErr bool
		err     error
	}{
		{
			name: "success",
			ids:  []int{4, 5},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.DeleteById(context.Background(), tc.ids)

			if tc.wantErr {
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	repo := NewRouteRepo(testDbInstance)

	testCases := []struct {
		name    string
		data    entities.Route
		newPos  int
		wantErr bool
		err     error
	}{
		{
			name: "success (previously deleted id)",
			data: entities.Route{
				RouteID:   4,
				RouteName: "after_delete_route",
				Load:      1000.0,
				CargoType: "cargo_type",
			},
		},
		{
			name: "success (already existing id)",
			data: entities.Route{
				RouteID:   2,
				RouteName: "already_existing_id",
				Load:      2000.0,
				CargoType: "cargo_type_2",
			},
			newPos: 7,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Register(context.Background(), tc.data)

			if tc.wantErr {
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				if tc.newPos != 0 { // need to update route_id because previous is already existing
					tc.data.RouteID = tc.newPos
				}
				foundInDB, _ := repo.GetById(context.Background(), tc.data.RouteID)
				require.Nil(t, err)
				require.Equal(t, tc.data.RouteID, foundInDB.RouteID)
				require.Equal(t, tc.data.RouteName, foundInDB.RouteName)
				require.InEpsilon(t, tc.data.Load, foundInDB.Load, eps)
				require.Equal(t, tc.data.CargoType, foundInDB.CargoType)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	repo := NewRouteRepo(testDbInstance)

	testCases := []struct {
		name     string
		id       int
		expected entities.Route
		wantErr  bool
		err      error
	}{
		{
			name: "success",
			id:   6,
			expected: entities.Route{
				RouteID:   6,
				RouteName: "test6",
				Load:      6.0,
				CargoType: "cargo6",
			},
		},
		{
			name:    "not found",
			id:      5,
			wantErr: true,
			err:     fmt.Errorf("getting route by id: no rows in result set"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			route, err := repo.GetById(context.Background(), tc.id)

			if tc.wantErr {
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected.RouteID, route.RouteID)
				require.Equal(t, tc.expected.RouteName, route.RouteName)
				require.InEpsilon(t, tc.expected.Load, route.Load, eps)
				require.Equal(t, tc.expected.CargoType, route.CargoType)
			}
		})
	}
}
