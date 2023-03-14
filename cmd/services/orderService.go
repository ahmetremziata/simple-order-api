package services

import (
	"net/http"
	"simple-order-api/cmd/constants"
	enum "simple-order-api/cmd/enums"
	"simple-order-api/cmd/models/response"
	"simple-order-api/cmd/repositories"
)

type OrderService interface {
	GetOrders() ([]response.Order, *response.ErrorResponse)
	GetOrder(orderNumber string) (*response.Order, *response.ErrorResponse)
	DeleteOrder(orderNumber string) *response.ErrorResponse
}

type OrderServiceImp struct {
	orderRepository repositories.OrderRepository
}

func (o OrderServiceImp) GetOrders() ([]response.Order, *response.ErrorResponse) {
	return o.orderRepository.FetchOrders()
}

func (o OrderServiceImp) GetOrder(orderNumber string) (*response.Order, *response.ErrorResponse) {
	order, _ := o.orderRepository.FetchOrderByOrderNumber(orderNumber)
	return order, nil
}

func (o OrderServiceImp) DeleteOrder(orderNumber string) *response.ErrorResponse {
	order, errorResp := o.GetOrder(orderNumber)
	if errorResp != nil {
		return errorResp
	}

	if order == nil {
		errorResp := response.NewErrorBuilder().
			SetError(http.StatusNotFound, constants.OrderNotFoundByOrderNumber).
			Build()
		return &errorResp
	}

	if order.StatusId == int(enum.Transferred) ||
		order.StatusId == int(enum.Shipped) ||
		order.StatusId == int(enum.Delivered) {
		errorResp := response.NewErrorBuilder().
			SetError(http.StatusInternalServerError, constants.OrderDeletionNotPermittedBecauseOfStatus).
			Build()
		return &errorResp
	}

	return o.orderRepository.DeleteOrder(orderNumber)
}

func NewOrderService(orderRepository repositories.OrderRepository) OrderService {
	return &OrderServiceImp{
		orderRepository: orderRepository,
	}
}
