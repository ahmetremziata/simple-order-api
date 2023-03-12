package services

import (
	"simple-order-api/cmd/models/response"
	"simple-order-api/cmd/repositories"
)

type OrderService interface {
	GetOrders() ([]response.Order, *error)
	GetOrder(orderNumber string) (*response.Order, *error)
}

type OrderServiceImp struct {
	orderRepository repositories.OrderRepository
}

func (o OrderServiceImp) GetOrder(orderNumber string) (*response.Order, *error) {
	order, _ := o.orderRepository.FetchOrderByOrderNumber(orderNumber)
	return order, nil
}

func (o OrderServiceImp) GetOrders() ([]response.Order, *error) {
	return o.orderRepository.FetchOrders()
}

func NewOrderService(orderRepository repositories.OrderRepository) OrderService {
	return &OrderServiceImp{
		orderRepository: orderRepository,
	}
}
