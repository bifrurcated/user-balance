package sort

import (
	"context"
	"net/http"
	"strings"
)

const (
	ASC               = "ASC"
	DESC              = "DESC"
	OptionsContextKey = "sort_option"
)

func Middleware(h http.HandlerFunc, defaultSortField, defaultSortOrder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := r.URL.Query().Get("sort_order")
		if sortBy == "" {
			sortBy = defaultSortField
		}
		if sortOrder == "" {
			sortOrder = defaultSortOrder
		} else {
			sortOrderUpper := strings.ToUpper(sortOrder)
			if sortOrderUpper != ASC && sortOrderUpper != DESC {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("bad sort order"))
				// TODO w.Write - error for your api
				return
			}
		}
		options := Options{
			Field: sortBy,
			Order: sortOrder,
		}
		ctx := context.WithValue(r.Context(), OptionsContextKey, options)
		r = r.WithContext(ctx)

		h(w, r)
	}
}

type Options struct {
	Field, Order string
}
