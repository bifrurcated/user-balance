package reserve

import (
	"encoding/json"
	"github.com/bifrurcated/user-balance/internal/apperror"
	"github.com/bifrurcated/user-balance/internal/handlers"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	apiVersionURL = "/api/v1"
	reserveURL    = "/reserve"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handler {
	return &handler{service: service, logger: logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, apiVersionURL+reserveURL, apperror.Middleware(h.ReserveMoney))
}

func (h *handler) ReserveMoney(w http.ResponseWriter, r *http.Request) error {
	var reserveDTO CreateReserveDTO
	err := json.NewDecoder(r.Body).Decode(&reserveDTO)
	if err != nil {
		return apperror.NewAppError(nil, err.Error(), "", "US-000004")
	}
	reserve, err := h.service.ReserveMoney(r.Context(), reserveDTO)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(reserve)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) ReserveProfit(w http.ResponseWriter, r *http.Request) error {

	return nil
}
