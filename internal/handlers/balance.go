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

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawRequest struct {
	OrderID string  `json:"order"`
	Sum     float64 `json:"sum"`
}

type WithdrawalOrder struct {
	OrderID     string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type WithdrawalResponse []*WithdrawalOrder

func GetUsersBalance(s Storage) http.HandlerFunc {
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

		currentBalance, err := s.GetUsersCurrentBalance(r.Context(), userID)
		if err != nil {
			logger.Log.Error("Error getting user balance", zap.Error(err))
			http.Error(w, "Error getting user balance", http.StatusInternalServerError)
			return
		}
		withdrawnBalance, err := s.GetUsersWithdrawnBalance(r.Context(), userID)
		if err != nil {
			logger.Log.Error("Error getting user withdrawn balance", zap.Error(err))
			http.Error(w, "Error getting user withdrawn balance", http.StatusInternalServerError)
			return
		}

		resp := &BalanceResponse{
			Current:   currentBalance,
			Withdrawn: withdrawnBalance,
		}

		data, err := json.Marshal(&resp)
		if err != nil {
			logger.Log.Error("Error marshaling response", zap.Error(err))
			http.Error(w, "Error marshaling response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func WithdrawOrder(s Storage) http.HandlerFunc {
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

		var buf bytes.Buffer
		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			logger.Log.Error("Error reading body", zap.Error(err))
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		req := &WithdrawRequest{}
		err = json.Unmarshal(buf.Bytes(), &req)
		if err != nil {
			logger.Log.Error("Error unmarshalling body", zap.Error(err))
			http.Error(w, "Error unmarshalling body", http.StatusInternalServerError)
			return
		}

		if req.Sum <= 0 {
			http.Error(w, "Invalid sum", http.StatusBadRequest)
			return
		}

		currentBalance, err := s.GetUsersCurrentBalance(r.Context(), userID)
		if err != nil {
			logger.Log.Error("Error getting users current balance", zap.Error(err))
			http.Error(w, "Error getting users current balance", http.StatusInternalServerError)
			return
		}

		if currentBalance < req.Sum {
			http.Error(w, "Not enough balance", http.StatusPaymentRequired)
			return
		}

		orderID, err := strconv.Atoi(req.OrderID)
		if err != nil {
			logger.Log.Error("Error parsing order id", zap.Error(err))
			http.Error(w, "Error parsing order id", http.StatusInternalServerError)
			return
		}

		order := &models.Order{
			ID:         orderID,
			UserID:     userID,
			Status:     models.OrderStatusProcessed,
			Sum:        -req.Sum,
			UploadedAt: time.Now(),
		}
		err = order.Validate()
		if err != nil {
			http.Error(w, "Invalid order id", http.StatusUnprocessableEntity)
			return
		}

		_, err = s.CreateOrder(r.Context(), order)
		if err != nil {
			if errors.Is(err, storage.ErrOrderAlreadyExists) {
				http.Error(w, "Order with this number already exists", http.StatusConflict)
				return
			}
			if errors.Is(err, storage.ErrOrderCreatedByOtherUser) {
				http.Error(w, "Order was already created by other user", http.StatusConflict)
				return
			}
			logger.Log.Error("Error creating order", zap.Error(err))
			http.Error(w, "Error creating order", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetUsersWithdrawals(s Storage) http.HandlerFunc {
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

		withdrawals, err := s.GetUsersWithdrawals(r.Context(), userID)
		if err != nil {
			logger.Log.Error("Error getting user withdrawals", zap.Error(err))
			http.Error(w, "Error getting user withdrawals", http.StatusInternalServerError)
			return
		}

		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp := make(WithdrawalResponse, len(withdrawals))
		for i, v := range withdrawals {
			resp[i] = &WithdrawalOrder{
				OrderID:     strconv.Itoa(v.ID),
				Sum:         v.Sum,
				ProcessedAt: v.UploadedAt.Format(time.RFC3339),
			}
		}

		data, err := json.Marshal(&resp)
		if err != nil {
			logger.Log.Error("Error marshaling response", zap.Error(err))
			http.Error(w, "Error marshaling response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
