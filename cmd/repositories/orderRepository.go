package repositories

import (
	"simple-order-api/cmd/models/response"
)

type OrderRepository interface {
	FetchOrders() ([]response.Order, *response.ErrorResponse)
	FetchOrderByOrderNumber(orderNumber string) (*response.Order, *response.ErrorResponse)
	DeleteOrder(orderNumber string) *response.ErrorResponse
}

type OrderRepositoryImp struct{}

func (o OrderRepositoryImp) FetchOrders() ([]response.Order, *response.ErrorResponse) {
	orders := getOrders()
	return orders, nil
}

func (o OrderRepositoryImp) FetchOrderByOrderNumber(orderNumber string) (*response.Order, *response.ErrorResponse) {
	orders := getOrders()
	for _, order := range orders {
		if order.OrderNumber == orderNumber {
			return &order, nil
		}
	}

	return nil, nil
}

func (o OrderRepositoryImp) DeleteOrder(_ string) *response.ErrorResponse {
	// Think that we delete order with success
	return nil
}

// This function represents data for external service
func getOrders() []response.Order {
	orders := []response.Order{
		{
			OrderNumber:  "1",
			FirstName:    "Ahmet",
			LastName:     "Ata",
			TotalAmount:  121.13,
			Address:      "Lorem ipsum dolor sit amet",
			City:         "İstanbul",
			District:     "Silivri",
			StatusId:     2,
			CurrencyCode: "TR",
		},
		{
			OrderNumber:  "2",
			FirstName:    "Hans",
			LastName:     "Schengen",
			TotalAmount:  345.99,
			Address:      "Sed ut perspiciatis unde omnis iste natus",
			City:         "Berlin",
			District:     "Berlin Square",
			StatusId:     3,
			CurrencyCode: "EUR",
		},
		{
			OrderNumber:  "3",
			FirstName:    "George",
			LastName:     "White",
			TotalAmount:  163.99,
			Address:      "Ut enim ad minima veniam, quis nostrum",
			City:         "London",
			District:     "Birmingham",
			StatusId:     4,
			CurrencyCode: "EUR",
		},
	}
	return orders
}

func NewOrderRepository() OrderRepository {
	return &OrderRepositoryImp{}
}
