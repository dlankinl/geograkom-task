package services

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"task/internal/dto"
	"task/internal/entities"
	"task/internal/mocks"
	"testing"
	"time"
)

func TestDeleteByIds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRouteRepo(ctrl)
	svc := NewRouteService(repo)

	testCases := []struct {
		beforeTest func(repo mocks.MockRouteRepo)
		name       string
		ids        dto.DeleteRoutesRequestBody
		wantErr    bool
		err        error
	}{
		{
			name:    "success",
			wantErr: false,
			ids:     dto.DeleteRoutesRequestBody{RouteIDs: []int{1, 2, 3}},
			beforeTest: func(repo mocks.MockRouteRepo) {
				repo.EXPECT().DeleteById(gomock.Any(), []int{1, 2, 3}).Return(nil)
			},
		},
		{
			name:    "empty ids",
			wantErr: false,
			ids:     dto.DeleteRoutesRequestBody{RouteIDs: []int{}},
			beforeTest: func(repo mocks.MockRouteRepo) {
				repo.EXPECT().DeleteById(gomock.Any(), []int{}).Return(nil)
			},
		},
		{
			name:    "negative id",
			wantErr: true,
			ids:     dto.DeleteRoutesRequestBody{RouteIDs: []int{1, -2, 3}},
			err:     fmt.Errorf("deleting routes: ids should be non-negative"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*repo)
			}

			err := svc.DeleteByIds(context.Background(), tc.ids)

			if tc.wantErr {
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}

			// goroutine needs time to complete ;)
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestGetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRouteRepo(ctrl)
	svc := NewRouteService(repo)

	testCases := []struct {
		name       string
		id         int
		expected   entities.Route
		beforeTest func(repo mocks.MockRouteRepo)
		wantErr    bool
		err        error
	}{
		{
			name: "success",
			id:   1,
			expected: entities.Route{
				RouteID:   1,
				RouteName: "test",
				Load:      1000.0,
				CargoType: "sand",
			},
			beforeTest: func(repo mocks.MockRouteRepo) {
				repo.EXPECT().GetById(gomock.Any(), 1).Return(entities.Route{
					RouteID:   1,
					RouteName: "test",
					Load:      1000.0,
					CargoType: "sand",
					IsActual:  false,
				}, nil)
			},
		},
		{
			name:    "id is negative",
			id:      -1,
			wantErr: true,
			err:     fmt.Errorf("route id should be non-negative"),
		},
		{
			name: "error in repository",
			id:   1,
			beforeTest: func(repo mocks.MockRouteRepo) {
				repo.EXPECT().
					GetById(gomock.Any(), 1).
					Return(
						entities.Route{},
						fmt.Errorf("some repo error"),
					)
			},
			wantErr: true,
			err:     fmt.Errorf("getting route by id: some repo error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*repo)
			}

			route, err := svc.GetById(context.Background(), tc.id)

			if tc.wantErr {
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.Equal(t, tc.expected, route)
				require.Nil(t, err)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRouteRepo(ctrl)
	svc := NewRouteService(repo)

	testCases := []struct {
		name       string
		data       dto.RegisterRouteRequestBody
		beforeTest func(repo mocks.MockRouteRepo)
		wantErr    bool
		err        error
	}{
		{
			name: "success",
			data: dto.RegisterRouteRequestBody{
				RouteID:   1,
				RouteName: "test",
				Load:      1000.0,
				CargoType: "sand",
			},
			beforeTest: func(repo mocks.MockRouteRepo) {
				repo.EXPECT().
					Register(
						gomock.Any(),
						entities.Route{
							RouteID:   1,
							RouteName: "test",
							Load:      1000.0,
							CargoType: "sand",
						}).
					Return(false, nil)
			},
		},
		{
			name: "negative load",
			data: dto.RegisterRouteRequestBody{
				RouteID:   1,
				RouteName: "test",
				Load:      -1000.0,
				CargoType: "sand",
			},
			wantErr: true,
			err:     fmt.Errorf("converting dto to entity model: load should be non-negative"),
		},
		{
			name: "error in repository",
			data: dto.RegisterRouteRequestBody{
				RouteID:   1,
				RouteName: "test",
				Load:      1000.0,
				CargoType: "sand",
			},
			beforeTest: func(repo mocks.MockRouteRepo) {
				repo.EXPECT().
					Register(
						gomock.Any(),
						entities.Route{
							RouteID:   1,
							RouteName: "test",
							Load:      1000.0,
							CargoType: "sand",
						}).
					Return(false, fmt.Errorf("some repo error"))
			},
			wantErr: true,
			err:     fmt.Errorf("route registration: some repo error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*repo)
			}

			_, err := svc.Register(context.Background(), tc.data)

			if tc.wantErr {
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
