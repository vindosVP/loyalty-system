package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/vindosVP/loyalty-system/internal/models"
	"github.com/vindosVP/loyalty-system/internal/storage"
	"github.com/vindosVP/loyalty-system/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type OrderResponse struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	Withdrawal float64 `json:"withdrawal,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type OrdersListResponse []*OrderResponse

func CreateOrder(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			logger.Log.Error("Error reading body", zap.Error(err))
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		gotOrderID := buf.String()
		if len(gotOrderID) == 0 {
			http.Error(w, "Empty order id", http.StatusBadRequest)
			return
		}
		orderID, err := strconv.Atoi(gotOrderID)
		if err != nil {
			http.Error(w, "Invalid order id", http.StatusBadRequest)
			return
		}

		gotUserID := r.Header.Get("x-user-id")
		if gotUserID == "" {
			logger.Log.Error("User id is empty")
			http.Error(w, "User id is empty", http.StatusInternalServerError)
			return
		}
		userID, err := strconv.Atoi(gotUserID)
		if err != nil {
			logger.Log.Error("Error parsing user id", zap.Error(err))
			http.Error(w, "Error parsing user id", http.StatusInternalServerError)
			return
		}

		order := &models.Order{
			ID:         orderID,
			UserID:     userID,
			Status:     models.OrderStatusNew,
			Sum:        0,
			UploadedAt: time.Now(),
		}

		err = order.Validate()
		if err != nil {
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}

		_, err = s.CreateOrder(r.Context(), order)
		if err != nil {
			if errors.Is(err, storage.ErrOrderAlreadyExists) {
				w.WriteHeader(http.StatusOK)
				return
			}
			if errors.Is(err, storage.ErrOrderCreatedByOtherUser) {
				http.Error(w, "Order was already created by other user", http.StatusConflict)
				return
			}
			logger.Log.Error("Error creating order", zap.Error(err))
			http.Error(w, "Error creating order", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func GetOrderList(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		gotUserID := r.Header.Get("x-user-id")
		if gotUserID == "" {
			logger.Log.Error("User id is empty")
			http.Error(w, "User id is empty", http.StatusInternalServerError)
			return
		}
		userID, err := strconv.Atoi(gotUserID)
		if err != nil {
			logger.Log.Error("Error parsing user id", zap.Error(err))
			http.Error(w, "Error parsing user id", http.StatusInternalServerError)
			return
		}

		usersOrders, err := s.GetUsersOrders(r.Context(), userID)
		if err != nil {
			logger.Log.Error("Error getting users orders", zap.Error(err))
			http.Error(w, "Error getting users orders", http.StatusInternalServerError)
			return
		}

		if len(usersOrders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp := make([]*OrderResponse, len(usersOrders))
		for i, order := range usersOrders {
			resp[i] = &OrderResponse{
				Number:     strconv.Itoa(order.ID),
				Status:     order.Status,
				UploadedAt: order.UploadedAt.Format(time.RFC3339),
			}
			if order.Sum > 0 {
				resp[i].Accrual = order.Sum
			}
			if order.Sum < 0 {
				resp[i].Withdrawal = -order.Sum
			}
		}

		data, err := json.Marshal(&resp)
		if err != nil {
			logger.Log.Error("Error marshaling orders", zap.Error(err))
			http.Error(w, "Error marshaling orders", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			logger.Log.Error("Error writing orders", zap.Error(err))
			http.Error(w, "Error writing orders", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
