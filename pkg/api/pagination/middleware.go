package pagination

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
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
			amount, err := strconv.ParseUint(last, 10, 64)
			if err != nil {
				parseTime, err2 := time.Parse(time.RFC3339, last)
				if err2 != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error() + " AND " + err2.Error()))
					return
				}
				value = parseTime
			} else {
				value = amount
			}
			fmt.Println(value)
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
