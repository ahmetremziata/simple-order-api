package mocks

import (
	"github.com/stretchr/testify/mock"
	"simple-order-api/cmd/models/response"
)

type FakeOrderService struct {
	mock.Mock
}

func (service *FakeOrderService) GetOrders() ([]response.Order, *error) {
	result := service.Called()
	if result.Get(0) != nil {
		return result.Get(0).([]response.Order), nil
	}

	return nil, result.Get(1).(*error)
}

func (service *FakeOrderService) GetOrder(orderNumber string) (*response.Order, *error) {
	result := service.Called(orderNumber)
	if result.Get(0) != nil {
		return result.Get(0).(*response.Order), nil
	}

	return nil, result.Get(1).(*error)
}
