package model

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidOrderInput = errors.New("invalid order input")
	ErrOrderAlreadyPaid  = errors.New("order already paid")
	ErrOrderCancelled    = errors.New("order is cancelled")
)
