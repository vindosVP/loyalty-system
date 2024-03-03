package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/vindosVP/loyalty-system/internal/models"
	"github.com/vindosVP/loyalty-system/internal/storage"
	"github.com/vindosVP/loyalty-system/pkg/logger"
	"github.com/vindosVP/loyalty-system/pkg/passwords"
	"github.com/vindosVP/loyalty-system/pkg/tokens"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Storage interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}

type LoginResponse struct {
	AuthToken string `json:"authToken"`
}

func Login(s Storage, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			logger.Log.Error("Error reading body", zap.Error(err))
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		user := &models.User{}
		err = json.Unmarshal(buf.Bytes(), &user)
		if err != nil {
			logger.Log.Error("Error unmarshalling body", zap.Error(err))
			http.Error(w, "Error unmarshalling body", http.StatusInternalServerError)
			return
		}

		if err = user.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		gotUser, err := s.GetUserByLogin(r.Context(), user.Login)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				http.Error(w, "Invalid login or password", http.StatusUnauthorized)
				return
			}
			logger.Log.Error("Error getting user", zap.Error(err))
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}
		if !passwords.Compare(user.Pwd, gotUser.EncryptedPwd) {
			http.Error(w, "Invalid login or password", http.StatusUnauthorized)
			return
		}

		token, err := tokens.CreateJwt(
			tokens.JWTClaims(gotUser.Id, gotUser.Login, time.Now().Add(time.Hour*72).Unix()),
			jwtSecret,
		)
		if err != nil {
			logger.Log.Error("Error creating token", zap.Error(err))
			http.Error(w, "Error creating token", http.StatusInternalServerError)
			return
		}

		resp := &LoginResponse{
			AuthToken: token,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			logger.Log.Error("Error marshalling response", zap.Error(err))
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
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

func Create(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			logger.Log.Error("Error reading body", zap.Error(err))
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		user := &models.User{}
		err = json.Unmarshal(buf.Bytes(), &user)
		if err != nil {
			logger.Log.Error("Error unmarshalling body", zap.Error(err))
			http.Error(w, "Error unmarshalling body", http.StatusInternalServerError)
			return
		}

		if err = user.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		encPwd, err := passwords.Encrypt(user.Pwd)
		if err != nil {
			logger.Log.Error("Error encrypting password", zap.Error(err))
			http.Error(w, "Error encrypting password", http.StatusInternalServerError)
			return
		}

		user.EncryptedPwd = encPwd

		createdUser, err := s.CreateUser(r.Context(), user)
		if err != nil {
			if errors.Is(err, storage.ErrUserAlreadyExists) {
				http.Error(w, "User already exists", http.StatusConflict)
				return
			}
			logger.Log.Error("Error creating user", zap.Error(err))
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(createdUser)
		if err != nil {
			logger.Log.Error("Error marshalling user", zap.Error(err))
			http.Error(w, "Error marshalling user", http.StatusInternalServerError)
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
