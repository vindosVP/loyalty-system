package storage

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vindosVP/loyalty-system/internal/models"
	"github.com/vindosVP/loyalty-system/internal/storage/mocks"
	"testing"
	"time"
)

func TestStorage_CreateUser(t *testing.T) {
	unexpectedError := errors.New("unexpected error")

	type userRepoExistsMock struct {
		needed bool
		result bool
		err    error
	}
	type userRepoCreateMock struct {
		needed bool
		result *models.User
		err    error
	}
	type args struct {
		user *models.User
	}
	type want struct {
		result *models.User
		err    error
	}

	tests := []struct {
		name               string
		userRepoExistsMock userRepoExistsMock
		userRepoCreateMock userRepoCreateMock
		args               args
		want               want
	}{
		{
			name: "ok",
			userRepoExistsMock: userRepoExistsMock{
				needed: true,
				result: false,
				err:    nil,
			},
			userRepoCreateMock: userRepoCreateMock{
				needed: true,
				result: &models.User{
					ID:           1,
					Login:        "testUser",
					EncryptedPwd: "encryptedPwd",
				},
				err: nil,
			},
			args: args{
				user: &models.User{
					ID:           1,
					Login:        "testUser",
					EncryptedPwd: "encryptedPwd",
				},
			},
			want: want{
				result: &models.User{
					ID:           1,
					Login:        "testUser",
					EncryptedPwd: "encryptedPwd",
				},
				err: nil,
			},
		},
		{
			name: "user already exists",
			userRepoExistsMock: userRepoExistsMock{
				needed: true,
				result: true,
				err:    nil,
			},
			userRepoCreateMock: userRepoCreateMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			args: args{
				user: &models.User{
					ID:           1,
					Login:        "testUser",
					EncryptedPwd: "encryptedPwd",
				},
			},
			want: want{
				result: nil,
				err:    ErrUserAlreadyExists,
			},
		},
		{
			name: "user exists unexpected error",
			userRepoExistsMock: userRepoExistsMock{
				needed: true,
				result: false,
				err:    unexpectedError,
			},
			userRepoCreateMock: userRepoCreateMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			args: args{
				user: &models.User{
					ID:           1,
					Login:        "testUser",
					EncryptedPwd: "encryptedPwd",
				},
			},
			want: want{
				result: nil,
				err:    unexpectedError,
			},
		},
		{
			name: "creation unexpected error",
			userRepoExistsMock: userRepoExistsMock{
				needed: true,
				result: false,
				err:    nil,
			},
			userRepoCreateMock: userRepoCreateMock{
				needed: true,
				result: nil,
				err:    unexpectedError,
			},
			args: args{
				user: &models.User{
					ID:           1,
					Login:        "testUser",
					EncryptedPwd: "encryptedPwd",
				},
			},
			want: want{
				result: nil,
				err:    unexpectedError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo := mocks.NewUserRepo(t)
			orderRepo := mocks.NewOrderRepo(t)
			s := New(userRepo, orderRepo)

			if tt.userRepoExistsMock.needed {
				userRepo.On("Exists", mock.Anything, tt.args.user.Login).Return(tt.userRepoExistsMock.result, tt.userRepoExistsMock.err)
			}
			if tt.userRepoCreateMock.needed {
				userRepo.On("Create", mock.Anything, tt.args.user).Return(tt.userRepoCreateMock.result, tt.userRepoCreateMock.err)
			}

			result, err := s.CreateUser(ctx, tt.args.user)

			if tt.want.result != nil {
				assert.Equal(t, tt.want.result, result)
				assert.NoError(t, err)
			}

			if tt.want.err != nil {
				assert.ErrorIs(t, err, tt.want.err)
			}
		})
	}
}

func TestStorage_GetUserByLogin(t *testing.T) {
	unexpectedError := errors.New("unexpected error")

	type userRepoExistsMock struct {
		needed bool
		result bool
		err    error
	}
	type userRepoGetByLoginMock struct {
		needed bool
		result *models.User
		err    error
	}
	type args struct {
		login string
	}
	type want struct {
		result *models.User
		err    error
	}

	tests := []struct {
		name                   string
		userRepoExistsMock     userRepoExistsMock
		userRepoGetByLoginMock userRepoGetByLoginMock
		args                   args
		want                   want
	}{
		{
			name: "ok",
			userRepoExistsMock: userRepoExistsMock{
				needed: true,
				result: true,
				err:    nil,
			},
			userRepoGetByLoginMock: userRepoGetByLoginMock{
				needed: true,
				result: &models.User{
					ID:           1,
					Login:        "testUser",
					EncryptedPwd: "encryptedPwd",
				},
				err: nil,
			},
			args: args{
				login: "testUser",
			},
			want: want{
				result: &models.User{
					ID:           1,
					Login:        "testUser",
					EncryptedPwd: "encryptedPwd",
				},
				err: nil,
			},
		},
		{
			name: "user not found",
			userRepoExistsMock: userRepoExistsMock{
				needed: true,
				result: false,
				err:    nil,
			},
			userRepoGetByLoginMock: userRepoGetByLoginMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			args: args{
				login: "testUser",
			},
			want: want{
				result: nil,
				err:    ErrUserNotFound,
			},
		},
		{
			name: "user exists unexpected error",
			userRepoExistsMock: userRepoExistsMock{
				needed: true,
				result: false,
				err:    unexpectedError,
			},
			userRepoGetByLoginMock: userRepoGetByLoginMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			args: args{
				login: "testUser",
			},
			want: want{
				result: nil,
				err:    unexpectedError,
			},
		},
		{
			name: "userRepo.GetByLogin unexpected error",
			userRepoExistsMock: userRepoExistsMock{
				needed: true,
				result: true,
				err:    nil,
			},
			userRepoGetByLoginMock: userRepoGetByLoginMock{
				needed: true,
				result: nil,
				err:    unexpectedError,
			},
			args: args{
				login: "testUser",
			},
			want: want{
				result: nil,
				err:    unexpectedError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo := mocks.NewUserRepo(t)
			orderRepo := mocks.NewOrderRepo(t)
			s := New(userRepo, orderRepo)

			if tt.userRepoExistsMock.needed {
				userRepo.On("Exists", mock.Anything, tt.args.login).Return(tt.userRepoExistsMock.result, tt.userRepoExistsMock.err)
			}
			if tt.userRepoGetByLoginMock.needed {
				userRepo.On("GetByLogin", mock.Anything, tt.args.login).Return(tt.userRepoGetByLoginMock.result, tt.userRepoGetByLoginMock.err)
			}

			result, err := s.GetUserByLogin(ctx, tt.args.login)

			if tt.want.result != nil {
				assert.Equal(t, tt.want.result, result)
				assert.NoError(t, err)
			}

			if tt.want.err != nil {
				assert.ErrorIs(t, err, tt.want.err)
			}
		})
	}
}

func TestStorage_CreateOrder(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	currentTime := time.Now()

	type OrderRepoExistsMock struct {
		needed bool
		result bool
		err    error
	}
	type OrderRepoGetByIDMock struct {
		needed bool
		result *models.Order
		err    error
	}
	type OrderRepoCreateMock struct {
		needed bool
		result *models.Order
		err    error
	}
	type args struct {
		order *models.Order
	}
	type want struct {
		result *models.Order
		err    error
	}

	tests := []struct {
		name                 string
		orderRepoExistsMock  OrderRepoExistsMock
		orderRepoGetByIDMock OrderRepoGetByIDMock
		orderRepoCreateMock  OrderRepoCreateMock
		args                 args
		want                 want
	}{
		{
			name: "ok",
			orderRepoExistsMock: OrderRepoExistsMock{
				needed: true,
				result: false,
				err:    nil,
			},
			orderRepoGetByIDMock: OrderRepoGetByIDMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			orderRepoCreateMock: OrderRepoCreateMock{
				needed: true,
				result: &models.Order{
					ID:         1,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
				err: nil,
			},
			args: args{
				order: &models.Order{
					ID:         1,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
			},
			want: want{
				result: &models.Order{
					ID:         1,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
				err: nil,
			},
		},
		{
			name: "order already exists",
			orderRepoExistsMock: OrderRepoExistsMock{
				needed: true,
				result: true,
				err:    nil,
			},
			orderRepoGetByIDMock: OrderRepoGetByIDMock{
				needed: true,
				result: &models.Order{
					ID:         1,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
				err: nil,
			},
			orderRepoCreateMock: OrderRepoCreateMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			args: args{
				order: &models.Order{
					ID:         1,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
			},
			want: want{
				result: nil,
				err:    ErrOrderAlreadyExists,
			},
		},
		{
			name: "order already created by other user",
			orderRepoExistsMock: OrderRepoExistsMock{
				needed: true,
				result: true,
				err:    nil,
			},
			orderRepoGetByIDMock: OrderRepoGetByIDMock{
				needed: true,
				result: &models.Order{
					ID:         1,
					UserID:     2,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
				err: nil,
			},
			orderRepoCreateMock: OrderRepoCreateMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			args: args{
				order: &models.Order{
					ID:         1,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
			},
			want: want{
				result: nil,
				err:    ErrOrderCreatedByOtherUser,
			},
		},
		{
			name: "orderRepo.Exists unexpected error",
			orderRepoExistsMock: OrderRepoExistsMock{
				needed: true,
				result: false,
				err:    unexpectedError,
			},
			orderRepoGetByIDMock: OrderRepoGetByIDMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			orderRepoCreateMock: OrderRepoCreateMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			args: args{
				order: &models.Order{
					ID:         1,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
			},
			want: want{
				result: nil,
				err:    unexpectedError,
			},
		},
		{
			name: "orderRepo.GetByID unexpected error",
			orderRepoExistsMock: OrderRepoExistsMock{
				needed: true,
				result: true,
				err:    nil,
			},
			orderRepoGetByIDMock: OrderRepoGetByIDMock{
				needed: true,
				result: nil,
				err:    unexpectedError,
			},
			orderRepoCreateMock: OrderRepoCreateMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			args: args{
				order: &models.Order{
					ID:         1,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
			},
			want: want{
				result: nil,
				err:    unexpectedError,
			},
		},
		{
			name: "orderRepo.Create unexpected error",
			orderRepoExistsMock: OrderRepoExistsMock{
				needed: true,
				result: false,
				err:    nil,
			},
			orderRepoGetByIDMock: OrderRepoGetByIDMock{
				needed: false,
				result: nil,
				err:    nil,
			},
			orderRepoCreateMock: OrderRepoCreateMock{
				needed: true,
				result: nil,
				err:    unexpectedError,
			},
			args: args{
				order: &models.Order{
					ID:         1,
					UserID:     1,
					Status:     models.OrderStatusNew,
					Sum:        0,
					UploadedAt: currentTime,
				},
			},
			want: want{
				result: nil,
				err:    unexpectedError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo := mocks.NewUserRepo(t)
			orderRepo := mocks.NewOrderRepo(t)
			s := New(userRepo, orderRepo)
			if tt.orderRepoExistsMock.needed {
				orderRepo.On("Exists", mock.Anything, tt.args.order.ID).Return(tt.orderRepoExistsMock.result, tt.orderRepoExistsMock.err)
			}
			if tt.orderRepoGetByIDMock.needed {
				orderRepo.On("GetByID", mock.Anything, tt.args.order.ID).Return(tt.orderRepoGetByIDMock.result, tt.orderRepoGetByIDMock.err)
			}
			if tt.orderRepoCreateMock.needed {
				orderRepo.On("Create", mock.Anything, tt.args.order).Return(tt.orderRepoCreateMock.result, tt.orderRepoCreateMock.err)
			}

			result, err := s.CreateOrder(ctx, tt.args.order)
			if tt.want.result != nil {
				assert.Equal(t, tt.want.result, result)
				assert.NoError(t, err)
			}
			if tt.want.err != nil {
				assert.ErrorIs(t, err, tt.want.err)
			}
		})
	}
}

func TestStorage_GetUsersOrders(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	currentTime := time.Now()

	type orderRepoGetUsersOrdersMock struct {
		needed bool
		result []*models.Order
		err    error
	}
	type args struct {
		userID int
	}
	type want struct {
		result []*models.Order
		err    error
	}

	tests := []struct {
		name                        string
		orderRepoGetUsersOrdersMock orderRepoGetUsersOrdersMock
		args                        args
		want                        want
	}{
		{
			name: "ok",
			orderRepoGetUsersOrdersMock: orderRepoGetUsersOrdersMock{
				needed: true,
				result: []*models.Order{
					{
						ID:         1,
						UserID:     1,
						Status:     models.OrderStatusNew,
						Sum:        0,
						UploadedAt: currentTime,
					},
					{
						ID:         2,
						UserID:     1,
						Status:     models.OrderStatusNew,
						Sum:        0,
						UploadedAt: currentTime,
					},
				},
				err: nil,
			},
			args: args{
				userID: 1,
			},
			want: want{
				result: []*models.Order{
					{
						ID:         1,
						UserID:     1,
						Status:     models.OrderStatusNew,
						Sum:        0,
						UploadedAt: currentTime,
					},
					{
						ID:         2,
						UserID:     1,
						Status:     models.OrderStatusNew,
						Sum:        0,
						UploadedAt: currentTime,
					},
				},
				err: nil,
			},
		},
		{
			name: "orderRepo.GetUsersOrders unexpected error",
			orderRepoGetUsersOrdersMock: orderRepoGetUsersOrdersMock{
				needed: true,
				result: nil,
				err:    unexpectedError,
			},
			args: args{
				userID: 1,
			},
			want: want{
				result: nil,
				err:    unexpectedError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo := mocks.NewUserRepo(t)
			orderRepo := mocks.NewOrderRepo(t)
			s := New(userRepo, orderRepo)

			if tt.orderRepoGetUsersOrdersMock.needed {
				orderRepo.On("GetUsersOrders", mock.Anything, tt.args.userID).Return(tt.orderRepoGetUsersOrdersMock.result, tt.orderRepoGetUsersOrdersMock.err)
			}

			result, err := s.GetUsersOrders(ctx, tt.args.userID)
			if tt.want.result != nil {
				assert.Equal(t, tt.want.result, result)
				assert.NoError(t, err)
			}
			if tt.want.err != nil {
				assert.ErrorIs(t, err, tt.want.err)
			}
		})
	}
}

func TestStorage_GetUsersCurrentBalance(t *testing.T) {
	unexpectedError := errors.New("unexpected error")

	type orderRepoGetUsersCurrentBalanceMock struct {
		needed bool
		result float64
		err    error
	}
	type args struct {
		userID int
	}
	type want struct {
		result float64
		err    error
	}

	tests := []struct {
		name                                string
		orderRepoGetUsersCurrentBalanceMock orderRepoGetUsersCurrentBalanceMock
		args                                args
		want                                want
	}{
		{
			name: "ok",
			orderRepoGetUsersCurrentBalanceMock: orderRepoGetUsersCurrentBalanceMock{
				needed: true,
				result: 512.3,
				err:    nil,
			},
			args: args{
				userID: 1,
			},
			want: want{
				result: 512.3,
				err:    nil,
			},
		},
		{
			name: "orderRepo.GetUsersCurrentBalance unexpected error",
			orderRepoGetUsersCurrentBalanceMock: orderRepoGetUsersCurrentBalanceMock{
				needed: true,
				result: 0,
				err:    unexpectedError,
			},
			args: args{
				userID: 1,
			},
			want: want{
				result: 0,
				err:    unexpectedError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo := mocks.NewUserRepo(t)
			orderRepo := mocks.NewOrderRepo(t)
			s := New(userRepo, orderRepo)

			if tt.orderRepoGetUsersCurrentBalanceMock.needed {
				orderRepo.On("GetUsersCurrentBalance", mock.Anything, tt.args.userID).Return(tt.orderRepoGetUsersCurrentBalanceMock.result, tt.orderRepoGetUsersCurrentBalanceMock.err)
			}

			result, err := s.GetUsersCurrentBalance(ctx, tt.args.userID)
			if tt.want.err == nil {
				assert.Equal(t, tt.want.result, result)
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.want.err)
			}
		})
	}
}

func TestStorage_GetUsersWithdrawnBalance(t *testing.T) {
	unexpectedError := errors.New("unexpected error")

	type orderRepoGetUsersWithdrawnBalanceMock struct {
		needed bool
		result float64
		err    error
	}
	type args struct {
		userID int
	}
	type want struct {
		result float64
		err    error
	}

	tests := []struct {
		name                                  string
		orderRepoGetUsersWithdrawnBalanceMock orderRepoGetUsersWithdrawnBalanceMock
		args                                  args
		want                                  want
	}{
		{
			name: "ok",
			orderRepoGetUsersWithdrawnBalanceMock: orderRepoGetUsersWithdrawnBalanceMock{
				needed: true,
				result: 512.3,
				err:    nil,
			},
			args: args{
				userID: 1,
			},
			want: want{
				result: 512.3,
				err:    nil,
			},
		},
		{
			name: "orderRepo.GetUsersWithdrawnBalance unexpected error",
			orderRepoGetUsersWithdrawnBalanceMock: orderRepoGetUsersWithdrawnBalanceMock{
				needed: true,
				result: 0,
				err:    unexpectedError,
			},
			args: args{
				userID: 1,
			},
			want: want{
				result: 0,
				err:    unexpectedError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo := mocks.NewUserRepo(t)
			orderRepo := mocks.NewOrderRepo(t)
			s := New(userRepo, orderRepo)

			if tt.orderRepoGetUsersWithdrawnBalanceMock.needed {
				orderRepo.On("GetUsersWithdrawnBalance", mock.Anything, tt.args.userID).Return(tt.orderRepoGetUsersWithdrawnBalanceMock.result, tt.orderRepoGetUsersWithdrawnBalanceMock.err)
			}

			result, err := s.GetUsersWithdrawnBalance(ctx, tt.args.userID)
			if tt.want.err == nil {
				assert.Equal(t, tt.want.result, result)
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.want.err)
			}
		})
	}
}

func TestStorage_GetUsersWithdrawals(t *testing.T) {
	unexpectedError := errors.New("unexpected error")
	currentTime := time.Now()

	type orderRepoGetUsersWithdrawalsMock struct {
		needed bool
		result []*models.Order
		err    error
	}
	type args struct {
		userID int
	}
	type want struct {
		result []*models.Order
		err    error
	}

	tests := []struct {
		name                             string
		orderRepoGetUsersWithdrawalsMock orderRepoGetUsersWithdrawalsMock
		args                             args
		want                             want
	}{
		{
			name: "ok",
			orderRepoGetUsersWithdrawalsMock: orderRepoGetUsersWithdrawalsMock{
				needed: true,
				result: []*models.Order{
					{
						ID:         1,
						UserID:     1,
						Status:     models.OrderStatusProcessed,
						Sum:        100,
						UploadedAt: currentTime,
					},
					{
						ID:         2,
						UserID:     1,
						Status:     models.OrderStatusProcessed,
						Sum:        200,
						UploadedAt: currentTime,
					},
				},
				err: nil,
			},
			args: args{
				userID: 1,
			},
			want: want{
				result: []*models.Order{
					{
						ID:         1,
						UserID:     1,
						Status:     models.OrderStatusProcessed,
						Sum:        100,
						UploadedAt: currentTime,
					},
					{
						ID:         2,
						UserID:     1,
						Status:     models.OrderStatusProcessed,
						Sum:        200,
						UploadedAt: currentTime,
					},
				},
				err: nil,
			},
		},
		{
			name: "orderRepo.GetUsersWithdrawals unexpected error",
			orderRepoGetUsersWithdrawalsMock: orderRepoGetUsersWithdrawalsMock{
				needed: true,
				result: nil,
				err:    unexpectedError,
			},
			args: args{
				userID: 1,
			},
			want: want{
				result: nil,
				err:    unexpectedError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo := mocks.NewUserRepo(t)
			orderRepo := mocks.NewOrderRepo(t)
			s := New(userRepo, orderRepo)

			if tt.orderRepoGetUsersWithdrawalsMock.needed {
				orderRepo.On("GetUsersWithdrawals", mock.Anything, tt.args.userID).Return(tt.orderRepoGetUsersWithdrawalsMock.result, tt.orderRepoGetUsersWithdrawalsMock.err)
			}

			result, err := s.GetUsersWithdrawals(ctx, tt.args.userID)
			if tt.want.result != nil {
				assert.Equal(t, tt.want.result, result)
				assert.NoError(t, err)
			}
			if tt.want.err != nil {
				assert.ErrorIs(t, err, tt.want.err)
			}
		})
	}
}
