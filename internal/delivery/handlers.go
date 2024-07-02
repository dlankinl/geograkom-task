package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task/internal/app"
	"task/internal/dto"
)

func RegisterHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "register handler"

		var req dto.RegisterRouteRequestBody

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		err = app.Svc.Register(r.Context(), req)
		if err != nil {
			errorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		// TODO: 208 code
		successResponse(w, http.StatusOK, nil)
	}
}

func GetHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "get handler"

		id := r.URL.Query().Get("id")
		if id == "" {
			errorResponse(w, fmt.Errorf("%s: empty id", prompt).Error(), http.StatusInternalServerError)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("%s: converting string id to int: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		route, err := app.Svc.GetById(r.Context(), idInt)
		if err != nil {
			errorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		if !route.IsActual {
			errorResponse(w, fmt.Errorf("%s: route is not actual", prompt).Error(), http.StatusGone)
			return
		}

		successResponse(w, http.StatusOK, map[string]any{
			"route_name": route.RouteName,
			"load":       route.Load,
			"cargo_type": route.CargoType,
		})
	}
}

func DeleteHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "delete handler"

		var req dto.DeleteRoutesRequestBody

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		err = app.Svc.DeleteByIds(r.Context(), req)
		if err != nil {
			errorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusAccepted, nil)
	}
}
