package balance

import (
	"encoding/json"
	"fmt"
	"github.com/bifrurcated/user-balance/internal/apperror"
	"github.com/bifrurcated/user-balance/internal/handlers"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	apiVersionURL = "/api/v1"
	balanceURL    = "/balance"
	addMoneyURL   = "/add-money"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handler {
	return &handler{service: service, logger: logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, apiVersionURL+balanceURL, apperror.Middleware(h.GetBalance))
	router.HandlerFunc(http.MethodPost, apiVersionURL+addMoneyURL, apperror.Middleware(h.AddMoney))
}

func (h *handler) AddMoney(w http.ResponseWriter, r *http.Request) error {
	var tum TransferUserMoney
	err := json.NewDecoder(r.Body).Decode(&tum)
	if err != nil {
		return apperror.NewAppError(nil, err.Error(), "", "US-000004")
	}
	h.logger.Debugf("tum = %v", tum)
	if tum.Amount < 0 {
		return apperror.NewAppError(
			nil,
			fmt.Sprintf("cannot add negative amount (%f)", tum.Amount),
			"", "US-000004")
	}
	err = h.service.AddAmount(r.Context(), tum)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *handler) GetBalance(w http.ResponseWriter, r *http.Request) error {
	var ub UserBalance
	err := json.NewDecoder(r.Body).Decode(&ub)
	if err != nil {
		return apperror.NewAppError(nil, err.Error(), "", "US-000004")
	}
	one, err := h.service.GetOne(r.Context(), ub.UserID)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(one.Amount)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return nil
}
