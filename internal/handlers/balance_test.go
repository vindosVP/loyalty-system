package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vindosVP/loyalty-system/internal/handlers/mocks"
	"github.com/vindosVP/loyalty-system/internal/models"
	"github.com/vindosVP/loyalty-system/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetUsersBalance(t *testing.T) {
	uri := "/api/user/balance"

	type request struct {
		method string
		userID string
	}
	type getUsersCurrentBalanceMock struct {
		needed bool
		result float64
		err    error
	}
	type getUsersWithdrawnBalanceMock struct {
		needed bool
		result float64
		err    error
	}
	type want struct {
		statusCode int
		result     BalanceResponse
	}

	tests := []struct {
		name                         string
		request                      request
		getUsersCurrentBalanceMock   getUsersCurrentBalanceMock
		getUsersWithdrawnBalanceMock getUsersWithdrawnBalanceMock
		want                         want
	}{
		{
			name: "ok",
			request: request{
				method: http.MethodGet,
				userID: "1",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: true,
				result: 100,
				err:    nil,
			},
			getUsersWithdrawnBalanceMock: getUsersWithdrawnBalanceMock{
				needed: true,
				result: 50,
				err:    nil,
			},
			want: want{
				statusCode: http.StatusOK,
				result: BalanceResponse{
					Current:   100,
					Withdrawn: 50,
				},
			},
		},
		{
			name: "wrong method",
			request: request{
				method: http.MethodPost,
				userID: "1",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: false,
				result: 0,
				err:    nil,
			},
			getUsersWithdrawnBalanceMock: getUsersWithdrawnBalanceMock{
				needed: false,
				result: 0,
				err:    nil,
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				result: BalanceResponse{
					Current:   0,
					Withdrawn: 0,
				},
			},
		},
		{
			name: "zero balance",
			request: request{
				method: http.MethodGet,
				userID: "1",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: true,
				result: 0,
				err:    nil,
			},
			getUsersWithdrawnBalanceMock: getUsersWithdrawnBalanceMock{
				needed: true,
				result: 0,
				err:    nil,
			},
			want: want{
				statusCode: http.StatusOK,
				result: BalanceResponse{
					Current:   0,
					Withdrawn: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewStorage(t)
			if tt.getUsersCurrentBalanceMock.needed {
				s.On("GetUsersCurrentBalance", mock.Anything, mock.Anything).Return(tt.getUsersCurrentBalanceMock.result, tt.getUsersCurrentBalanceMock.err)
			}
			if tt.getUsersWithdrawnBalanceMock.needed {
				s.On("GetUsersWithdrawnBalance", mock.Anything, mock.Anything).Return(tt.getUsersWithdrawnBalanceMock.result, tt.getUsersWithdrawnBalanceMock.err)
			}

			r := chi.NewRouter()
			r.Get("/api/user/balance", GetUsersBalance(s))

			req := httptest.NewRequest(tt.request.method, uri, nil)
			req.Header.Set("x-user-id", tt.request.userID)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			if tt.want.statusCode == http.StatusOK {
				var response BalanceResponse
				err := json.NewDecoder(res.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.result, response)
			}
		})
	}
}

func TestWithdrawOrder(t *testing.T) {
	uri := "/api/user/withdraw"

	type request struct {
		method string
		userID string
		body   string
	}
	type getUsersCurrentBalanceMock struct {
		needed bool
		result float64
		err    error
	}
	type createOrderMock struct {
		needed bool
		result *models.Order
		err    error
	}
	type want struct {
		statusCode int
	}

	tests := []struct {
		name                       string
		request                    request
		getUsersCurrentBalanceMock getUsersCurrentBalanceMock
		createOrderMock            createOrderMock
		want                       want
	}{
		{
			name: "ok",
			request: request{
				method: http.MethodPost,
				userID: "1",
				body:   "{\"order\": \"8023459525\", \"sum\": 100}",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: true,
				result: 200,
				err:    nil,
			},
			createOrderMock: createOrderMock{
				needed: true,
				result: &models.Order{
					ID:         1,
					UserID:     8023459525,
					Status:     models.OrderStatusProcessed,
					Sum:        -100,
					UploadedAt: time.Now(),
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "wrong method",
			request: request{
				method: http.MethodGet,
				userID: "1",
				body:   "{\"order\": \"8023459525\", \"sum\": 100}",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: false,
				result: 0,
				err:    nil,
			},
			createOrderMock: createOrderMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "wrong order number",
			request: request{
				method: http.MethodPost,
				userID: "1",
				body:   "{\"order\": \"1111\", \"sum\": 100}",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: true,
				result: 200,
				err:    nil,
			},
			createOrderMock: createOrderMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "wrong sum",
			request: request{
				method: http.MethodPost,
				userID: "1",
				body:   "{\"order\": \"8023459525\", \"sum\": -100}",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: false,
				result: 0,
				err:    nil,
			},
			createOrderMock: createOrderMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "not enough balance",
			request: request{
				method: http.MethodPost,
				userID: "1",
				body:   "{\"order\": \"8023459525\", \"sum\": 100}",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: true,
				result: 50,
				err:    nil,
			},
			createOrderMock: createOrderMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			want: want{
				statusCode: http.StatusPaymentRequired,
			},
		},
		{
			name: "order already exists",
			request: request{
				method: http.MethodPost,
				userID: "1",
				body:   "{\"order\": \"8023459525\", \"sum\": 100}",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: true,
				result: 200,
				err:    nil,
			},
			createOrderMock: createOrderMock{
				needed: true,
				result: nil,
				err:    storage.ErrOrderAlreadyExists,
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		{
			name: "order already created by other user",
			request: request{
				method: http.MethodPost,
				userID: "1",
				body:   "{\"order\": \"8023459525\", \"sum\": 100}",
			},
			getUsersCurrentBalanceMock: getUsersCurrentBalanceMock{
				needed: true,
				result: 200,
				err:    nil,
			},
			createOrderMock: createOrderMock{
				needed: true,
				result: nil,
				err:    storage.ErrOrderCreatedByOtherUser,
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewStorage(t)
			if tt.getUsersCurrentBalanceMock.needed {
				s.On("GetUsersCurrentBalance", mock.Anything, mock.Anything).Return(tt.getUsersCurrentBalanceMock.result, tt.getUsersCurrentBalanceMock.err)
			}
			if tt.createOrderMock.needed {
				s.On("CreateOrder", mock.Anything, mock.Anything).Return(tt.createOrderMock.result, tt.createOrderMock.err)
			}

			r := chi.NewRouter()
			r.Post(uri, WithdrawOrder(s))
			req := httptest.NewRequest(tt.request.method, uri, strings.NewReader(tt.request.body))
			req.Header.Set("x-user-id", tt.request.userID)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}

func TestGetUsersWithdrawals(t *testing.T) {
	uri := "/api/user/withdrawals"
	currentTime := time.Now()

	type request struct {
		method string
		userID string
	}
	type want struct {
		statusCode int
		result     WithdrawalResponse
	}
	type getUsersWithdrawalsMock struct {
		needed bool
		result []*models.Order
		err    error
	}

	tests := []struct {
		name                    string
		request                 request
		getUsersWithdrawalsMock getUsersWithdrawalsMock
		want                    want
	}{
		{
			name: "ok",
			request: request{
				method: http.MethodGet,
				userID: "1",
			},
			getUsersWithdrawalsMock: getUsersWithdrawalsMock{
				needed: true,
				result: []*models.Order{
					{
						ID:         1,
						UserID:     1,
						Status:     models.OrderStatusProcessed,
						Sum:        200,
						UploadedAt: currentTime,
					},
					{
						ID:         2,
						UserID:     1,
						Status:     models.OrderStatusProcessed,
						Sum:        1,
						UploadedAt: currentTime,
					}, {
						ID:         3,
						UserID:     1,
						Status:     models.OrderStatusProcessed,
						Sum:        1000,
						UploadedAt: currentTime,
					},
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				result: WithdrawalResponse{
					&WithdrawalOrder{
						OrderID:     "1",
						Sum:         200,
						ProcessedAt: currentTime.Format(time.RFC3339),
					},
					&WithdrawalOrder{
						OrderID:     "2",
						Sum:         1,
						ProcessedAt: currentTime.Format(time.RFC3339),
					},
					&WithdrawalOrder{
						OrderID:     "3",
						Sum:         1000,
						ProcessedAt: currentTime.Format(time.RFC3339),
					},
				},
			},
		},
		{
			name: "no withdrawals",
			request: request{
				method: http.MethodGet,
				userID: "1",
			},
			getUsersWithdrawalsMock: getUsersWithdrawalsMock{
				needed: true,
				result: make([]*models.Order, 0),
				err:    nil,
			},
			want: want{
				statusCode: http.StatusNoContent,
				result:     nil,
			},
		},
		{
			name: "wrong method",
			request: request{
				method: http.MethodPost,
				userID: "1",
			},
			getUsersWithdrawalsMock: getUsersWithdrawalsMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				result:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewStorage(t)
			if tt.getUsersWithdrawalsMock.needed {
				s.On("GetUsersWithdrawals", mock.Anything, mock.Anything).Return(tt.getUsersWithdrawalsMock.result, tt.getUsersWithdrawalsMock.err)
			}

			r := chi.NewRouter()
			r.Get(uri, GetUsersWithdrawals(s))
			req := httptest.NewRequest(tt.request.method, uri, nil)
			req.Header.Set("x-user-id", tt.request.userID)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			if tt.want.statusCode == http.StatusOK {
				var response WithdrawalResponse
				err := json.NewDecoder(res.Body).Decode(&response)
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.want.result, response)
			}
		})
	}
}
