package pagination

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"net/http"
	"strconv"
	"time"
)

const (
	OptionsContextKey = "pagination_option"
)

type Options struct {
	Limit uint64
	Value any
	ID    uint64
}

type Token struct {
	Value string `json:"value"`
	Type  string `json:"type"`
	ID    uint64 `json:"id"`
}

func Middleware(h http.HandlerFunc, defaultLimit uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger()
		logger.Info("start pagination middleware")
		limitStr := r.URL.Query().Get("limit")
		token := r.URL.Query().Get("token")
		limit := defaultLimit
		if limitStr != "" {
			parseLimit, err := strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			limit = parseLimit
		}
		var id uint64
		var value any
		if token != "" {
			decodeString, err := hex.DecodeString(token)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			logger.Trace(string(decodeString))
			t := Token{}
			err = json.Unmarshal(decodeString, &t)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			logger.Trace(t)
			if t.Type == "amount" {
				value, err = strconv.ParseFloat(t.Value, 32)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
			} else if t.Type == "datetime" {
				parseInt, parseErr := strconv.ParseInt(t.Value, 10, 64)
				if parseErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(parseErr.Error()))
					return
				}
				value = time.Unix(parseInt, 0)
			} else {
				value = t.Value
			}
			id = t.ID
		}
		options := Options{
			Limit: limit,
			Value: value,
			ID:    id,
		}
		logger.Trace(options)
		logger.Info("end pagination middleware")
		ctx := context.WithValue(r.Context(), OptionsContextKey, options)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
