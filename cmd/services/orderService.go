package services

import (
	"simple-order-api/cmd/models/response"
)

type OrderService interface {
	GetOrders() ([]response.Order, *error)
}

type OrderServiceImp struct{}

func (o OrderServiceImp) GetOrders() ([]response.Order, *error) {
	data := []response.Order{
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
			District:     "Little Berlin",
			StatusId:     3,
			CurrencyCode: "EUR",
		},
	}
	return data, nil
}

func NewOrderService() OrderService {
	return &OrderServiceImp{}
}
