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

func TestCreateOrder(t *testing.T) {
	uri := "/api/user/orders"

	type request struct {
		method string
		body   string
		userID string
	}
	type want struct {
		statusCode int
	}
	type createOrderMock struct {
		needed bool
		result *models.Order
		err    error
	}

	tests := []struct {
		name            string
		createOrderMock createOrderMock
		request         request
		want            want
	}{
		{
			name: "ok",
			createOrderMock: createOrderMock{
				needed: true,
				result: &models.Order{
					ID:         7703824164,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: time.Now(),
				},
				err: nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "7703824164",
				userID: "1",
			},
			want: want{
				statusCode: http.StatusAccepted,
			},
		},
		{
			name: "wrong method",
			createOrderMock: createOrderMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodGet,
				body:   "7703824164",
				userID: "1",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "invalid order number",
			createOrderMock: createOrderMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "1111",
				userID: "1",
			},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "empty order number",
			createOrderMock: createOrderMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "",
				userID: "1",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "order already exists",
			createOrderMock: createOrderMock{
				needed: true,
				result: nil,
				err:    storage.ErrOrderAlreadyExists,
			},
			request: request{
				method: http.MethodPost,
				body:   "7703824164",
				userID: "1",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "order already created by another user",
			createOrderMock: createOrderMock{
				needed: true,
				result: nil,
				err:    storage.ErrOrderCreatedByOtherUser,
			},
			request: request{
				method: http.MethodPost,
				body:   "7703824164",
				userID: "1",
			},
			want: want{
				statusCode: http.StatusConflict,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewStorage(t)
			if tt.createOrderMock.needed {
				s.On("CreateOrder", mock.Anything, mock.Anything).Return(tt.createOrderMock.result, tt.createOrderMock.err)
			}

			r := chi.NewRouter()
			r.Post(uri, CreateOrder(s))

			req := httptest.NewRequest(tt.request.method, uri, strings.NewReader(tt.request.body))
			req.Header.Set("x-user-id", tt.request.userID)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}

func TestGetOrderList(t *testing.T) {
	uri := "/api/users/orders"
	currentTime := time.Now()
	currentTimeStr := currentTime.Format(time.RFC3339)

	type request struct {
		method string
		userID string
	}
	type getUserOrdersMock struct {
		needed bool
		result []*models.Order
		err    error
	}
	type want struct {
		statusCode int
		result     OrdersListResponse
	}

	tests := []struct {
		name              string
		getUserOrdersMock getUserOrdersMock
		request           request
		want              want
	}{
		{
			name: "ok",
			getUserOrdersMock: getUserOrdersMock{
				needed: true,
				result: []*models.Order{
					{
						ID:         9278923470,
						UserID:     1,
						Status:     models.OrderStatusProcessed,
						Sum:        500,
						UploadedAt: currentTime,
					},
					{
						ID:         12345678903,
						UserID:     1,
						Status:     models.OrderStatusProcessing,
						Sum:        0,
						UploadedAt: currentTime,
					},
					{
						ID:         346436439,
						UserID:     1,
						Status:     models.OrderStatusInvalid,
						Sum:        0,
						UploadedAt: currentTime,
					},
				},
				err: nil,
			},
			request: request{
				method: http.MethodGet,
				userID: "1",
			},
			want: want{
				statusCode: http.StatusOK,
				result: OrdersListResponse{
					&OrderResponse{
						Number:     "9278923470",
						Status:     "PROCESSED",
						Accrual:    500,
						UploadedAt: currentTimeStr,
					},
					&OrderResponse{
						Number:     "12345678903",
						Status:     "PROCESSING",
						UploadedAt: currentTimeStr,
					},
					&OrderResponse{
						Number:     "346436439",
						Status:     "INVALID",
						UploadedAt: currentTimeStr,
					},
				},
			},
		},
		{
			name: "wrong method",
			getUserOrdersMock: getUserOrdersMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodPost,
				userID: "1",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				result:     nil,
			},
		},
		{
			name: "no orders",
			getUserOrdersMock: getUserOrdersMock{
				needed: true,
				result: make([]*models.Order, 0),
				err:    nil,
			},
			request: request{
				method: http.MethodGet,
				userID: "1",
			},
			want: want{
				statusCode: http.StatusNoContent,
				result:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewStorage(t)
			if tt.getUserOrdersMock.needed {
				s.On("GetUsersOrders", mock.Anything, mock.Anything).Return(tt.getUserOrdersMock.result, tt.getUserOrdersMock.err)
			}

			r := chi.NewRouter()
			r.Get(uri, GetOrderList(s))

			req := httptest.NewRequest(tt.request.method, uri, nil)
			req.Header.Set("x-user-id", tt.request.userID)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			if tt.want.statusCode == http.StatusOK {
				var ordersListResponse OrdersListResponse
				err := json.NewDecoder(res.Body).Decode(&ordersListResponse)
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.want.result, ordersListResponse)
			}
		})
	}
}
