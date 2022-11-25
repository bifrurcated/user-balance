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
	profitURL     = "/profit"
	cancelURL     = "/cancel"
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
	router.HandlerFunc(http.MethodPost, apiVersionURL+reserveURL+profitURL, apperror.Middleware(h.ReserveProfit))
	router.HandlerFunc(http.MethodPost, apiVersionURL+reserveURL+cancelURL, apperror.Middleware(h.CancelReserve))
}

func (h *handler) ReserveMoney(w http.ResponseWriter, r *http.Request) error {
	var reserveDTO CreateReserveDTO
	err := json.NewDecoder(r.Body).Decode(&reserveDTO)
	if err != nil {
		return apperror.NewAppError(nil, err.Error(), "", "US-000004")
	}
	reserve, err := h.service.ReserveMoney(r.Context(), &reserveDTO)
	if err != nil {
		return err
	}
	reserveMoneyDTO := CreateReserveMoneyDTO{
		UserID:    reserve.UserID,
		ServiceID: reserve.ServiceID,
		OrderID:   reserve.OrderID,
		Cost:      reserve.Cost,
	}
	bytes, err := json.Marshal(reserveMoneyDTO)
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
	var profit ProfitReserveDTO
	err := json.NewDecoder(r.Body).Decode(&profit)
	if err != nil {
		return apperror.NewAppError(nil, err.Error(), "", "US-000004")
	}
	reserve, err := h.service.SetProfit(r.Context(), &CreateReserveDTO{
		UserID:    profit.UserID,
		ServiceID: profit.ServiceID,
		OrderID:   profit.OrderID,
		IsProfit:  true,
	})
	if err != nil {
		return err
	}
	profit.ID = reserve.ID
	profit.Cost = reserve.Cost
	bytes, err := json.Marshal(profit)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (h *handler) CancelReserve(w http.ResponseWriter, r *http.Request) error {
	var reserveDTO CancelReserveDTO
	err := json.NewDecoder(r.Body).Decode(&reserveDTO)
	if err != nil {
		return apperror.NewAppError(nil, err.Error(), "", "US-000004")
	}
	err = h.service.CancelReserve(r.Context(), &reserveDTO)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
