package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vindosVP/loyalty-system/pkg/tokens"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthenticator_WithAuth(t *testing.T) {
	type AuthTestResponse struct {
		UserID string
	}
	JWTSecret := "superSecret"
	uri := "/testAuth"

	type user struct {
		ID    int
		Login string
	}
	type auth struct {
		addHeader bool
		schema    string
	}
	type want struct {
		code   int
		result AuthTestResponse
	}

	tests := []struct {
		name string
		auth auth
		user user
		want want
	}{
		{
			name: "ok",
			auth: auth{
				addHeader: true,
				schema:    "Bearer",
			},
			user: user{
				ID:    1,
				Login: "someLogin",
			},
			want: want{
				code: http.StatusOK,
				result: AuthTestResponse{
					UserID: "1",
				},
			},
		},
		{
			name: "no auth header",
			auth: auth{
				addHeader: false,
				schema:    "Bearer",
			},
			user: user{
				ID:    1,
				Login: "someLogin",
			},
			want: want{
				code: http.StatusUnauthorized,
				result: AuthTestResponse{
					UserID: "1",
				},
			},
		},
		{
			name: "wrong schema",
			auth: auth{
				addHeader: false,
				schema:    "Basic",
			},
			user: user{
				ID:    1,
				Login: "someLogin",
			},
			want: want{
				code: http.StatusUnauthorized,
				result: AuthTestResponse{
					UserID: "1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			handler := func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					userID := r.Header.Get("x-user-id")
					resp := AuthTestResponse{UserID: userID}
					data, _ := json.Marshal(&resp)
					_, _ = w.Write(data)
					w.WriteHeader(http.StatusOK)
				}
			}

			a := NewAuthenticator(JWTSecret)

			r := chi.NewRouter()
			r.Use(a.WithAuth)
			r.Get(uri, handler())

			req := httptest.NewRequest("GET", uri, nil)
			if tt.auth.addHeader {
				token, err := tokens.CreateJWT(
					tokens.JWTClaims(tt.user.ID, tt.user.Login, time.Now().Add(time.Hour*72).Unix()), JWTSecret)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", tt.auth.schema, token))
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			if tt.want.code == http.StatusOK {
				var resp AuthTestResponse
				err := json.NewDecoder(res.Body).Decode(&resp)
				require.NoError(t, err)
				assert.Equal(t, tt.want.result, resp)
			}
		})
	}
}
