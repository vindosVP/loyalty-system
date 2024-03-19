package auth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestParseBearerToken(t *testing.T) {
	type args struct {
		r *http.Request
	}
	type want struct {
		token string
		err   error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ok",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"Authorization": []string{"Bearer someToken"},
					},
				},
			},
			want: want{
				token: "someToken",
				err:   nil,
			},
		},
		{
			name: "no auth header",
			args: args{
				r: &http.Request{
					Header: http.Header{},
				},
			},
			want: want{
				token: "",
				err:   ErrNoAuthHeader,
			},
		},
		{
			name: "invalid auth format",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"Authorization": []string{"Basic someToken"},
					},
				},
			},
			want: want{
				token: "",
				err:   ErrInvalidAuthFormat,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ParseBearerToken(tt.args.r)
			if tt.want.err == nil {
				assert.Equal(t, tt.want.token, token)
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.want.err)
			}
		})
	}
}
