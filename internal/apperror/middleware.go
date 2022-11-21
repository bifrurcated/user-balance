package apperror

import (
	"errors"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					_, err = w.Write(ErrNotFound.Marshal())
					if err != nil {
						logging.GetLogger().Error(err)
					}
					return
				}

				err = err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				_, err = w.Write(appErr.Marshal())
				if err != nil {
					logging.GetLogger().Error(err)
				}
				return
			}
			w.WriteHeader(http.StatusTeapot)
			_, err = w.Write(systemError(err).Marshal())
			if err != nil {
				logging.GetLogger().Error(err)
			}
		}
	}
}
