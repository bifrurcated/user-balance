package pagination

import (
	"context"
	"net/http"
	"strconv"
)

const (
	OptionsContextKey = "pagination_option"
)

func Middleware(h http.HandlerFunc, defaultLimit uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		last := r.URL.Query().Get("last")
		limit := defaultLimit
		if limitStr != "" {
			parseLimit, err := strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("cannot parse limit to uint64"))
				return
			}
			limit = parseLimit
		}
		var value any
		if last != "" {
			parsePageID, err := strconv.ParseUint(last, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("cannot parse page_id to uint64"))
				return
			}
			value = parsePageID
		}
		options := Options{
			Limit: limit,
			Value: value,
		}
		ctx := context.WithValue(r.Context(), OptionsContextKey, options)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

type Options struct {
	Limit uint64
	Value any
}
