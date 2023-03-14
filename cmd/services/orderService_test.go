package services

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"simple-order-api/cmd/constants"
	enum "simple-order-api/cmd/enums"
	"simple-order-api/cmd/mocks"
	"simple-order-api/cmd/models/response"
	"testing"
)

func TestGetOrders(t *testing.T) {
	//Given
	mockOrderRepository := &mocks.FakeOrderRepository{}
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
	}
	mockOrderRepository.On("FetchOrders").Return(orders, nil)
	service := NewOrderService(mockOrderRepository)

	//When
	resp, err := service.GetOrders()

	//Then
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, orders, resp)
	mockOrderRepository.AssertNumberOfCalls(t, "FetchOrders", 1)
	mockOrderRepository.AssertCalled(t, "FetchOrders")
}

func TestGetOrders_WhenOrderRepositoryReturnsError_ReturnsError(t *testing.T) {
	//Given
	mockOrderRepository := &mocks.FakeOrderRepository{}
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()
	mockOrderRepository.On("FetchOrders").Return(nil, &serviceErr)
	service := NewOrderService(mockOrderRepository)

	//When
	resp, err := service.GetOrders()

	//Then
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, &serviceErr, err)
}

func TestGetOrder(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.FakeOrderRepository{}
	order := response.Order{}

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(&order, nil)
	service := NewOrderService(mockOrderRepository)

	//When
	resp, err := service.GetOrder(orderNumber)

	//Then
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, &order, resp)
	mockOrderRepository.AssertNumberOfCalls(t, "FetchOrderByOrderNumber", 1)
	mockOrderRepository.AssertCalled(t, "FetchOrderByOrderNumber", orderNumber)
}

func TestGetOrder_WhenOrderRepositoryReturnsError_ReturnsError(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.FakeOrderRepository{}
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()
	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(nil, &serviceErr)
	service := NewOrderService(mockOrderRepository)

	//When
	resp, err := service.GetOrder(orderNumber)

	//Then
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, &serviceErr, err)
}

func TestDeleteOrder(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.FakeOrderRepository{}
	order := response.Order{
		OrderNumber:  "1",
		FirstName:    "Ahmet",
		LastName:     "Ata",
		TotalAmount:  121.13,
		Address:      "Lorem ipsum dolor sit amet",
		City:         "İstanbul",
		District:     "Silivri",
		StatusId:     int(enum.Created),
		CurrencyCode: "TR",
	}

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(&order, nil)
	mockOrderRepository.On("DeleteOrder", mock.Anything).Return(nil)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.DeleteOrder(orderNumber)

	//Then
	assert.Nil(t, err)
	mockOrderRepository.AssertNumberOfCalls(t, "FetchOrderByOrderNumber", 1)
	mockOrderRepository.AssertCalled(t, "FetchOrderByOrderNumber", orderNumber)
	mockOrderRepository.AssertNumberOfCalls(t, "DeleteOrder", 1)
	mockOrderRepository.AssertCalled(t, "DeleteOrder", orderNumber)
}

func TestDeleteOrder_WhenOrderRepositoryGetMethodReturnsError_ReturnsError(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.FakeOrderRepository{}
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(nil, &serviceErr)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.DeleteOrder(orderNumber)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, &serviceErr, err)
}

func TestDeleteOrder_WhenOrderNotFoundInRepository_ReturnsError(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.FakeOrderRepository{}

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(nil, nil)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.DeleteOrder(orderNumber)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Equal(t, constants.OrderNotFoundByOrderNumber, err.Message)
}

func TestDeleteOrder_WhenStatusIsNotValidForDeletion_ReturnsError(t *testing.T) {
	for _, statusId := range statusIdsThatNotBeValidForDeletion {
		t.Run(fmt.Sprintf("StatusId:%d", statusId), func(t *testing.T) {
			orderNumber := "1"
			mockOrderRepository := &mocks.FakeOrderRepository{}
			order := response.Order{
				OrderNumber:  "1",
				FirstName:    "Ahmet",
				LastName:     "Ata",
				TotalAmount:  121.13,
				Address:      "Lorem ipsum dolor sit amet",
				City:         "İstanbul",
				District:     "Silivri",
				StatusId:     statusId,
				CurrencyCode: "TR",
			}

			mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(&order, nil)
			service := NewOrderService(mockOrderRepository)

			//When
			err := service.DeleteOrder(orderNumber)

			//Then
			assert.NotNil(t, err)
			assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
			assert.Equal(t, constants.OrderDeletionNotPermittedBecauseOfStatus, err.Message)
		})
	}
}

var statusIdsThatNotBeValidForDeletion = []int{
	int(enum.Transferred),
	int(enum.Shipped),
	int(enum.Delivered),
}

func TestDeleteOrder_WhenOrderRepositoryDeleteMethodReturnsError_ReturnsError(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.FakeOrderRepository{}
	order := response.Order{
		OrderNumber:  "1",
		FirstName:    "Ahmet",
		LastName:     "Ata",
		TotalAmount:  121.13,
		Address:      "Lorem ipsum dolor sit amet",
		City:         "İstanbul",
		District:     "Silivri",
		StatusId:     int(enum.Created),
		CurrencyCode: "TR",
	}

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(&order, nil)
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()
	mockOrderRepository.On("DeleteOrder", mock.Anything).Return(&serviceErr)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.DeleteOrder(orderNumber)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, &serviceErr, err)
}
