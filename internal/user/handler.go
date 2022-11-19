package user

import (
	"encoding/json"
	"fmt"
	"github.com/bifrurcated/user-balance/internal/apperror"
	"github.com/bifrurcated/user-balance/internal/handlers"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
)

const (
	usersURL       = "/users"
	userURL        = "/users/:id"
	userBalanceURL = "/users/:id/balance"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handler {
	return &handler{service: service, logger: logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, userBalanceURL, apperror.Middleware(h.GetBalance))
	router.HandlerFunc(http.MethodPost, userURL, apperror.Middleware(h.AddMoney))
}

func (h *handler) AddMoney(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *handler) GetBalance(w http.ResponseWriter, r *http.Request) error {
	split := strings.Split(r.URL.Path, "/")
	strID := split[2]
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		return apperror.NewAppError(nil, fmt.Sprintf("cannot parse %s to int64", strID), "", "US-000004")
	}
	one, err := h.service.GetOne(r.Context(), id)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(one.Amount)
	if err != nil {
		return err
	}
	w.WriteHeader(200)
	w.Write(bytes)
	return nil
}
