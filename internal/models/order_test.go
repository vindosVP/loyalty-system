package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOrder_Validate(t *testing.T) {
	type args struct {
		order *Order
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
				order: &Order{
					ID:         7324401889,
					UserID:     1,
					Status:     OrderStatusNew,
					Sum:        0,
					UploadedAt: time.Now(),
				},
			},
			want: want{
				valid: true,
			},
		},
		{
			name: "invalid",
			args: args{
				order: &Order{
					ID:         1111111111,
					UserID:     1,
					Status:     OrderStatusNew,
					Sum:        0,
					UploadedAt: time.Now(),
				},
			},
			want: want{
				valid: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.order.Validate()
			if tt.want.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
