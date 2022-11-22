package history

import (
	"encoding/json"
	"github.com/bifrurcated/user-balance/internal/apperror"
	"github.com/bifrurcated/user-balance/internal/handlers"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

const (
	apiVersionURL = "/api/v1"
	historyURL    = "/history"
)

func NewHandler(service *Service, logger *logging.Logger) handlers.Handler {
	return &handler{service: service, logger: logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, apiVersionURL+historyURL, apperror.Middleware(h.GetUserTransactions))
}

func (h *handler) GetUserTransactions(w http.ResponseWriter, r *http.Request) error {
	var userHistory UserHistoryDTO
	err := json.NewDecoder(r.Body).Decode(&userHistory)
	if err != nil {
		return apperror.NewAppError(nil, err.Error(), "", "US-000004")
	}
	histories, err := h.service.UserTransactions(r.Context(), userHistory.UserID)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(histories)
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
