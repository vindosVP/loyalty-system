package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	type args struct {
		user *User
	}
	type want struct {
		valid bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "valid",
			args: args{
				user: &User{
					ID:    1,
					Login: "someLogin",
					Pwd:   "somePassword",
				},
			},
			want: want{
				valid: true,
			},
		},
		{
			name: "invalid login",
			args: args{
				user: &User{
					ID:    1,
					Login: "",
					Pwd:   "somePassword",
				},
			},
			want: want{
				valid: false,
			},
		},
		{
			name: "invalid password",
			args: args{
				user: &User{
					ID:    1,
					Login: "someLogin",
					Pwd:   "",
				},
			},
			want: want{
				valid: false,
			},
		},
		{
			name: "invalid login and password",
			args: args{
				user: &User{
					ID:    1,
					Login: "",
					Pwd:   "",
				},
			},
			want: want{
				valid: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.user.Validate()
			if tt.want.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
