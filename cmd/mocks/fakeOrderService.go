package mocks

import (
	"github.com/stretchr/testify/mock"
	"simple-order-api/cmd/models/response"
)

type FakeOrderService struct {
	mock.Mock
}

func (service *FakeOrderService) GetOrders() ([]response.Order, *response.ErrorResponse) {
	result := service.Called()
	if result.Get(0) != nil {
		return result.Get(0).([]response.Order), nil
	}

	return nil, result.Get(1).(*response.ErrorResponse)
}

func (service *FakeOrderService) GetOrder(orderNumber string) (*response.Order, *response.ErrorResponse) {
	result := service.Called(orderNumber)
	if result.Get(0) != nil {
		return result.Get(0).(*response.Order), nil
	}

	return nil, result.Get(1).(*response.ErrorResponse)
}

func (service *FakeOrderService) DeleteOrder(orderNumber string) *response.ErrorResponse {
	result := service.Called(orderNumber)
	if result.Get(0) != nil {
		return result.Get(0).(*response.ErrorResponse)
	}

	return nil
}
