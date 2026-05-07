package v1

import (
	service "github.com.Saveliy12/microservices2.git/order/internal/service"
)

type API struct {
	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) *API {
	return &API{orderService: orderService}
}
