package handlers

import (
	"encoding/json"
	"github.com/vindosVP/loyalty-system/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

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
