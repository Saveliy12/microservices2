package model

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrInvalidOrderInput = errors.New("invalid order input")
var ErrOrderAlreadyPaid = errors.New("order already paid")
var ErrOrderCancelled = errors.New("order is cancelled")
