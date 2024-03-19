package storage

import "errors"

var (
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrUserNotFound            = errors.New("user not found")
	ErrOrderAlreadyExists      = errors.New("order already exists")
	ErrOrderCreatedByOtherUser = errors.New("order created by other user")
)
