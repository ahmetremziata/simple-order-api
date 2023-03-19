package repositories

import (
	enum "simple-order-api/cmd/enums"
	"simple-order-api/cmd/models/request"
	"simple-order-api/cmd/models/response"
)

//go:generate mockery --name=OrderRepository --structname=MockOrderRepository --output=../mocks --filename=fakeOrderRepositoryWithMockery.go
type OrderRepository interface {
	FetchOrders() ([]response.Order, *response.ErrorResponse)
	FetchOrderByOrderNumber(orderNumber string) (*response.Order, *response.ErrorResponse)
	CreateOrder(createOrderRequest request.CreateOrderRequest) *response.ErrorResponse
	UpdateOrder(orderNumber string, updateOrderRequest request.UpdateOrderRequest) *response.ErrorResponse
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

func (o OrderRepositoryImp) CreateOrder(createOrderRequest request.CreateOrderRequest) *response.ErrorResponse {
	orders, _ := o.FetchOrders()
	orders = append(orders, response.Order{
		OrderNumber:  createOrderRequest.OrderNumber,
		FirstName:    createOrderRequest.FirstName,
		LastName:     createOrderRequest.LastName,
		TotalAmount:  createOrderRequest.TotalAmount,
		Address:      createOrderRequest.Address,
		City:         createOrderRequest.City,
		District:     createOrderRequest.District,
		CurrencyCode: createOrderRequest.CurrencyCode,
		StatusId:     int(enum.Created),
	})
	return nil
}

func (o OrderRepositoryImp) UpdateOrder(orderNumber string, updateOrderRequest request.UpdateOrderRequest) *response.ErrorResponse {
	order, _ := o.FetchOrderByOrderNumber(orderNumber)
	order.City = updateOrderRequest.City
	order.District = updateOrderRequest.District
	order.Address = updateOrderRequest.Address
	order.CurrencyCode = updateOrderRequest.CurrencyCode
	order.TotalAmount = updateOrderRequest.TotalAmount
	order.FirstName = updateOrderRequest.FirstName
	order.LastName = updateOrderRequest.LastName
	return nil
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
