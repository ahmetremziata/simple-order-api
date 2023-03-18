package services

import (
	"net/http"
	"simple-order-api/cmd/constants"
	enum "simple-order-api/cmd/enums"
	"simple-order-api/cmd/models/request"
	"simple-order-api/cmd/models/response"
	"simple-order-api/cmd/repositories"
)

//go:generate mockery --name=OrderService --structname=MockOrderService --output=../mocks --filename=fakeOrderServiceWithMockery.go
type OrderService interface {
	GetOrders() ([]response.Order, *response.ErrorResponse)
	GetOrder(orderNumber string) (*response.Order, *response.ErrorResponse)
	DeleteOrder(orderNumber string) *response.ErrorResponse
	CreateOrder(createOrderRequest request.CreateOrderRequest) *response.ErrorResponse
}

type OrderServiceImp struct {
	orderRepository repositories.OrderRepository
}

func (o OrderServiceImp) GetOrders() ([]response.Order, *response.ErrorResponse) {
	return o.orderRepository.FetchOrders()
}

func (o OrderServiceImp) GetOrder(orderNumber string) (*response.Order, *response.ErrorResponse) {
	order, err := o.orderRepository.FetchOrderByOrderNumber(orderNumber)
	return order, err
}

func (o OrderServiceImp) CreateOrder(createOrderRequest request.CreateOrderRequest) *response.ErrorResponse {
	order, errorResp := o.GetOrder(createOrderRequest.OrderNumber)
	if errorResp != nil {
		return errorResp
	}

	if order != nil {
		errorResp := response.NewErrorBuilder().
			SetError(http.StatusConflict, constants.SameOrderFoundByUniqueId).
			Build()
		return &errorResp
	}

	return o.orderRepository.CreateOrder(createOrderRequest)
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

	deleteErr := o.orderRepository.DeleteOrder(orderNumber)
	return deleteErr
}

func NewOrderService(orderRepository repositories.OrderRepository) OrderService {
	return &OrderServiceImp{
		orderRepository: orderRepository,
	}
}
