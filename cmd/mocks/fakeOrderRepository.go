package mocks

import (
	"github.com/stretchr/testify/mock"
	"simple-order-api/cmd/models/request"
	"simple-order-api/cmd/models/response"
)

type FakeOrderRepository struct {
	mock.Mock
}

func (service *FakeOrderRepository) FetchOrders() ([]response.Order, *response.ErrorResponse) {
	result := service.Called()
	if result.Get(0) != nil {
		return result.Get(0).([]response.Order), nil
	}

	return nil, result.Get(1).(*response.ErrorResponse)
}

func (service *FakeOrderRepository) FetchOrderByOrderNumber(orderNumber string) (*response.Order, *response.ErrorResponse) {
	result := service.Called(orderNumber)
	if result.Get(0) == nil && result.Get(1) == nil {
		return nil, nil
	}

	if result.Get(0) != nil {
		return result.Get(0).(*response.Order), nil
	}

	return nil, result.Get(1).(*response.ErrorResponse)
}

func (service *FakeOrderRepository) CreateOrder(createOrderRequest request.CreateOrderRequest) *response.ErrorResponse {
	result := service.Called(createOrderRequest)
	if result.Get(0) != nil {
		return result.Get(0).(*response.ErrorResponse)
	}

	return nil
}

func (service *FakeOrderRepository) UpdateOrder(orderNumber string, updateOrderRequest request.UpdateOrderRequest) *response.ErrorResponse {
	result := service.Called(orderNumber, updateOrderRequest)
	if result.Get(0) != nil {
		return result.Get(0).(*response.ErrorResponse)
	}

	return nil
}

func (service *FakeOrderRepository) DeleteOrder(orderNumber string) *response.ErrorResponse {
	result := service.Called(orderNumber)
	if result.Get(0) != nil {
		return result.Get(0).(*response.ErrorResponse)
	}

	return nil
}
