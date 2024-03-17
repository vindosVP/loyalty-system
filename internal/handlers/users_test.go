package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vindosVP/loyalty-system/internal/handlers/mocks"
	"github.com/vindosVP/loyalty-system/internal/models"
	"github.com/vindosVP/loyalty-system/internal/storage"
	"github.com/vindosVP/loyalty-system/pkg/passwords"
	"github.com/vindosVP/loyalty-system/pkg/tokens"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	jwtSecret := "superSecret"
	uri := "/api/user/register"

	type request struct {
		method string
		body   string
	}
	type want struct {
		code     int
		checkJWT bool
		userID   string
	}
	type createUserMock struct {
		needed bool
		result *models.User
		err    error
	}
	type getUserByLoginMock struct {
		needed bool
		result *models.User
		err    error
	}

	tests := []struct {
		name               string
		createUserMock     createUserMock
		getUserByLoginMock getUserByLoginMock
		request            request
		want               want
	}{
		{
			name: "ok",
			createUserMock: createUserMock{
				needed: true,
				result: &models.User{
					ID:    1,
					Login: "someLogin",
					Pwd:   "somePassword",
				},
				err: nil,
			},
			getUserByLoginMock: getUserByLoginMock{
				needed: true,
				result: nil,
				err:    storage.ErrUserNotFound,
			},
			request: request{
				method: http.MethodPost,
				body:   "{\"login\": \"someLogin\",\"password\": \"somePassword\"}",
			},
			want: want{
				code:     http.StatusOK,
				checkJWT: true,
				userID:   "1",
			},
		},
		{
			name: "wrong method",
			createUserMock: createUserMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			getUserByLoginMock: getUserByLoginMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodGet,
				body:   "{\"login\": \"\",\"password\": \"somePassword\"}",
			},
			want: want{
				code:     http.StatusMethodNotAllowed,
				checkJWT: false,
				userID:   "",
			},
		},
		{
			name: "invalid login",
			createUserMock: createUserMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			getUserByLoginMock: getUserByLoginMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "{\"login\": \"\",\"password\": \"somePassword\"}",
			},
			want: want{
				code:     http.StatusBadRequest,
				checkJWT: false,
				userID:   "",
			},
		},
		{
			name: "invalid password",
			createUserMock: createUserMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			getUserByLoginMock: getUserByLoginMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "{\"login\": \"\",\"password\": \"somePassword\"}",
			},
			want: want{
				code:     http.StatusBadRequest,
				checkJWT: false,
				userID:   "",
			},
		},
		{
			name: "user already exists",
			createUserMock: createUserMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			getUserByLoginMock: getUserByLoginMock{
				needed: true,
				result: &models.User{
					ID:    1,
					Login: "someLogin",
					Pwd:   "somePassword",
				},
				err: nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "{\"login\": \"someLogin\",\"password\": \"somePassword\"}",
			},
			want: want{
				code:     http.StatusConflict,
				checkJWT: false,
				userID:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := mocks.NewStorage(t)

			if tt.getUserByLoginMock.needed {
				s.On("GetUserByLogin", mock.Anything, mock.Anything).Return(tt.getUserByLoginMock.result, tt.getUserByLoginMock.err)
			}
			if tt.createUserMock.needed {
				s.On("CreateUser", mock.Anything, mock.Anything).Return(tt.createUserMock.result, tt.createUserMock.err)
			}

			r := chi.NewRouter()
			r.Post(uri, Register(s, jwtSecret))

			req := httptest.NewRequest(tt.request.method, uri, strings.NewReader(tt.request.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			if tt.want.checkJWT {
				auth := res.Header.Get("Authorization")
				require.NotEmpty(t, auth)

				splitToken := strings.Split(auth, "Bearer ")
				require.Len(t, splitToken, 2)
				token := splitToken[1]

				valid, err := tokens.IsAuthorized(token, jwtSecret)
				require.NoError(t, err)
				assert.True(t, valid)

				id, err := tokens.ExtractID(token, jwtSecret)
				require.NoError(t, err)
				assert.Equal(t, tt.want.userID, id)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	jwtSecret := "superSecret"
	uri := "/api/user/login"
	encryptedSomePassword, _ := passwords.Encrypt("somePassword")

	type request struct {
		method string
		body   string
	}
	type want struct {
		code     int
		checkJWT bool
		userID   string
	}
	type getUserByLoginMock struct {
		needed bool
		result *models.User
		err    error
	}

	tests := []struct {
		name               string
		getUserByLoginMock getUserByLoginMock
		request            request
		want               want
	}{
		{
			name: "ok",
			getUserByLoginMock: getUserByLoginMock{
				needed: true,
				result: &models.User{
					ID:           1,
					Login:        "someLogin",
					EncryptedPwd: encryptedSomePassword,
				},
				err: nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "{\"login\": \"someLogin\",\"password\": \"somePassword\"}",
			},
			want: want{
				code:     http.StatusOK,
				checkJWT: true,
				userID:   "1",
			},
		},
		{
			name: "wrong method",
			getUserByLoginMock: getUserByLoginMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodGet,
				body:   "{\"login\": \"someLogin\",\"password\": \"somePassword\"}",
			},
			want: want{
				code:     http.StatusMethodNotAllowed,
				checkJWT: false,
				userID:   "",
			},
		},
		{
			name: "invalid login",
			getUserByLoginMock: getUserByLoginMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "{\"login\": \"\",\"password\": \"somePassword\"}",
			},
			want: want{
				code:     http.StatusBadRequest,
				checkJWT: false,
				userID:   "",
			},
		},
		{
			name: "invalid password",
			getUserByLoginMock: getUserByLoginMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "{\"login\": \"someLogin\",\"password\": \"\"}",
			},
			want: want{
				code:     http.StatusBadRequest,
				checkJWT: false,
				userID:   "",
			},
		},
		{
			name: "user not found",
			getUserByLoginMock: getUserByLoginMock{
				needed: true,
				result: nil,
				err:    storage.ErrUserNotFound,
			},
			request: request{
				method: http.MethodPost,
				body:   "{\"login\": \"someLogin\",\"password\": \"password\"}",
			},
			want: want{
				code:     http.StatusUnauthorized,
				checkJWT: false,
				userID:   "",
			},
		},
		{
			name: "wrong password",
			getUserByLoginMock: getUserByLoginMock{
				needed: true,
				result: &models.User{
					ID:           1,
					Login:        "someLogin",
					EncryptedPwd: "someWrongPassword",
				},
				err: nil,
			},
			request: request{
				method: http.MethodPost,
				body:   "{\"login\": \"someLogin\",\"password\": \"somePassword\"}",
			},
			want: want{
				code:     http.StatusUnauthorized,
				checkJWT: false,
				userID:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewStorage(t)
			if tt.getUserByLoginMock.needed {
				s.On("GetUserByLogin", mock.Anything, mock.Anything).Return(tt.getUserByLoginMock.result, tt.getUserByLoginMock.err)
			}

			r := chi.NewRouter()
			r.Post(uri, Login(s, jwtSecret))
			req := httptest.NewRequest(tt.request.method, uri, strings.NewReader(tt.request.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			if tt.want.checkJWT {
				auth := res.Header.Get("Authorization")
				require.NotEmpty(t, auth)

				splitToken := strings.Split(auth, "Bearer ")
				require.Len(t, splitToken, 2)
				token := splitToken[1]

				valid, err := tokens.IsAuthorized(token, jwtSecret)
				require.NoError(t, err)
				assert.True(t, valid)

				id, err := tokens.ExtractID(token, jwtSecret)
				require.NoError(t, err)
				assert.Equal(t, tt.want.userID, id)
			}
		})
	}
}
