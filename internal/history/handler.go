package history

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bifrurcated/user-balance/internal/apperror"
	"github.com/bifrurcated/user-balance/internal/handlers"
	"github.com/bifrurcated/user-balance/pkg/api/pagination"
	"github.com/bifrurcated/user-balance/pkg/api/sort"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
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
	router.HandlerFunc(http.MethodGet, apiVersionURL+historyURL,
		pagination.Middleware(
			sort.Middleware(apperror.Middleware(h.GetUserTransactions), "datetime", sort.ASC),
			10))
}

/*
GetUserTransactions

	Response example:
	{
		"histories":[
			{"id":1, "sender_user_id": 1, "user_id":4, "service_id": null, "amount":100, "type":"перевод", "datetime": "2022-11-23T15:04:01.572675Z"}
			{"id":2, "sender_user_id": null, "user_id":4, "service_id": null, "amount":333.99, "type":"зачисление", "datetime": "2022-11-23T15:04:28.906406Z"}
		],
		"next_page_id": 7b2276616c7565223a2231363639313139363533222c2274797065223a226461746574696d65227d
	}
*/
func (h *handler) GetUserTransactions(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("start get user transaction")
	var paginationOptions pagination.Options
	if options, ok := r.Context().Value(pagination.OptionsContextKey).(pagination.Options); ok {
		paginationOptions = options
	}
	var sortOptions sort.Options
	if options, ok := r.Context().Value(sort.OptionsContextKey).(sort.Options); ok {
		sortOptions = options
	}
	var userHistory UserHistoryDTO
	err := json.NewDecoder(r.Body).Decode(&userHistory)
	if err != nil {
		return apperror.NewAppError(nil, err.Error(), "", "US-000004")
	}
	histories, err := h.service.UserTransactions(r.Context(), userHistory.UserID, OptionsDTO{
		Limit: paginationOptions.Limit + 1,
		Value: paginationOptions.Value,
		ID:    paginationOptions.ID,
		Field: sortOptions.Field,
		Order: sortOptions.Order,
	})
	if err != nil {
		return err
	}
	var token pagination.Token
	var historiesDTO UserHistoriesDTO
	lengthArr := uint64(len(histories))
	h.logger.Debugf("pagination limit: %d", paginationOptions.Limit)
	h.logger.Debugf("histories length: %d", lengthArr)
	if lengthArr == paginationOptions.Limit+1 {
		typeOrder := strings.ToLower(sortOptions.Field)
		token.ID = histories[lengthArr-1].ID
		if typeOrder == "amount" {
			token.Value = strconv.FormatFloat(float64(histories[lengthArr-1].Amount), 'f', -1, 32)
			token.Type = "amount"
		} else if typeOrder == "datetime" {
			token.Value = strconv.FormatInt(histories[lengthArr-1].Datetime.Unix(), 10)
			token.Type = "datetime"
		} else {
			var historyMap map[string]History
			historyBytes, marshalErr := json.Marshal(histories[lengthArr-1])
			if marshalErr != nil {
				return marshalErr
			}
			unmarshalErr := json.Unmarshal(historyBytes, &historyMap)
			if unmarshalErr != nil {
				return unmarshalErr
			}

			token.Value = fmt.Sprint(historyMap[typeOrder])
			token.Type = typeOrder
		}
		h.logger.Trace(token)
		tokenBytes, errMarshal := json.Marshal(token)
		if errMarshal != nil {
			return errMarshal
		}
		encodeToString := hex.EncodeToString(tokenBytes)
		historiesDTO = UserHistoriesDTO{
			Histories: histories[:lengthArr-1],
			NextPage:  &encodeToString,
		}
	} else {
		historiesDTO.Histories = histories
	}
	bytes, err := json.Marshal(historiesDTO)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	h.logger.Info("end get user transaction")
	return nil
}
