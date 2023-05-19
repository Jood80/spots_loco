package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func ValidateQueryParams(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		lat := r.URL.Query().Get("lat")
		long := r.URL.Query().Get("long")
		radius := r.URL.Query().Get("radius")
		shape := r.URL.Query().Get("type")

		if lat == "" || long == "" || radius == "" || shape == "" {
			http.Error(w, "Missing query parameters", http.StatusBadRequest)
			return
		}

		latFlo, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			http.Error(w, "Invalid value for 'lat'", http.StatusBadRequest)
			return
		}

		longFlo, err := strconv.ParseFloat(long, 64)
		if err != nil {
			http.Error(w, "Invalid value for 'long'", http.StatusBadRequest)
			return
		}

		radiusFlo, err := strconv.ParseFloat(radius, 64)
		if err != nil {
			http.Error(w, "Invalid value for 'radius'", http.StatusBadRequest)
			return
		}

		if shape != "square" && shape != "circle" {
			http.Error(w, "Invalid value for 'type', it's either square or circle", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "lat", latFlo)
		ctx = context.WithValue(ctx, "long", longFlo)
		ctx = context.WithValue(ctx, "radius", radiusFlo)
		ctx = context.WithValue(ctx, "shape", shape)

		r = r.WithContext(ctx)

		next(w, r, ps)
	}
}
