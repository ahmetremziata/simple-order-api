package services

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"simple-order-api/cmd/constants"
	enum "simple-order-api/cmd/enums"
	"simple-order-api/cmd/mocks"
	"simple-order-api/cmd/models/request"
	"simple-order-api/cmd/models/response"
	"testing"
)

func TestGetOrders(t *testing.T) {
	//Given
	mockOrderRepository := &mocks.MockOrderRepository{}
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
	mockOrderRepository := &mocks.MockOrderRepository{}
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
	mockOrderRepository := &mocks.MockOrderRepository{}
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
	mockOrderRepository := &mocks.MockOrderRepository{}
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

func TestCreateOrder(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.MockOrderRepository{}
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(nil, nil)
	mockOrderRepository.On("CreateOrder", mock.Anything).Return(nil)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.CreateOrder(*serviceReq)

	//Then
	assert.Nil(t, err)
	mockOrderRepository.AssertNumberOfCalls(t, "FetchOrderByOrderNumber", 1)
	mockOrderRepository.AssertCalled(t, "FetchOrderByOrderNumber", orderNumber)
	mockOrderRepository.AssertNumberOfCalls(t, "CreateOrder", 1)
	mockOrderRepository.AssertCalled(t, "CreateOrder", *serviceReq)
}

func TestCreateOrder_WhenOrderRepositoryGetMethodReturnsError_ReturnsError(t *testing.T) {
	//Given
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	mockOrderRepository := &mocks.MockOrderRepository{}
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(nil, &serviceErr)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.CreateOrder(*serviceReq)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, &serviceErr, err)
}

func TestCreateOrder_WhenOrderAlreadyExistWithSameOrderNumber_ReturnsStatusConflict(t *testing.T) {
	//Given
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	mockOrderRepository := &mocks.MockOrderRepository{}
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
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.CreateOrder(*serviceReq)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.StatusCode)
	assert.Equal(t, constants.SameOrderFoundByUniqueId, err.Message)
}

func TestCreateOrder_WhenOrderRepositoryCreateMethodReturnsError_ReturnsError(t *testing.T) {
	//Given
	mockOrderRepository := &mocks.MockOrderRepository{}
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(nil, nil)
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()
	mockOrderRepository.On("CreateOrder", mock.Anything).Return(&serviceErr)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.CreateOrder(*serviceReq)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, &serviceErr, err)
}

func TestUpdateOrder(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.MockOrderRepository{}
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
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
	mockOrderRepository.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.UpdateOrder(orderNumber, *serviceReq)

	//Then
	assert.Nil(t, err)
	mockOrderRepository.AssertNumberOfCalls(t, "FetchOrderByOrderNumber", 1)
	mockOrderRepository.AssertCalled(t, "FetchOrderByOrderNumber", orderNumber)
	mockOrderRepository.AssertNumberOfCalls(t, "UpdateOrder", 1)
	mockOrderRepository.AssertCalled(t, "UpdateOrder", orderNumber, *serviceReq)
}

func TestUpdateOrder_WhenOrderRepositoryGetMethodReturnsError_ReturnsError(t *testing.T) {
	//Given
	orderNumber := "1"
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	mockOrderRepository := &mocks.MockOrderRepository{}
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(nil, &serviceErr)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.UpdateOrder(orderNumber, *serviceReq)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, &serviceErr, err)
}

func TestUpdateOrder_WhenOrderNotFoundInRepository_ReturnsNotFound(t *testing.T) {
	//Given
	orderNumber := "1"
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	mockOrderRepository := &mocks.MockOrderRepository{}

	mockOrderRepository.On("FetchOrderByOrderNumber", mock.Anything).Return(nil, nil)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.UpdateOrder(orderNumber, *serviceReq)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Equal(t, constants.OrderNotFoundByOrderNumber, err.Message)
}

func TestUpdateOrder_WhenStatusIsNotValidForDeletion_ReturnsError(t *testing.T) {
	for _, statusId := range statusIdsThatNotBeValidForChangable {
		t.Run(fmt.Sprintf("StatusId:%d", statusId), func(t *testing.T) {
			orderNumber := "1"
			serviceReq := &request.UpdateOrderRequest{}
			_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
			mockOrderRepository := &mocks.MockOrderRepository{}
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
			err := service.UpdateOrder(orderNumber, *serviceReq)

			//Then
			assert.NotNil(t, err)
			assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
			assert.Equal(t, constants.OrderChangeNotPermittedBecauseOfStatus, err.Message)
		})
	}
}

func TestUpdateOrder_WhenOrderRepositoryUpdateMethodReturnsError_ReturnsError(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.MockOrderRepository{}
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)

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
	mockOrderRepository.On("UpdateOrder", mock.Anything, mock.Anything).Return(&serviceErr)
	service := NewOrderService(mockOrderRepository)

	//When
	err := service.UpdateOrder(orderNumber, *serviceReq)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, &serviceErr, err)
}

func TestDeleteOrder(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.MockOrderRepository{}
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
	mockOrderRepository := &mocks.MockOrderRepository{}
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
	mockOrderRepository := &mocks.MockOrderRepository{}

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
			mockOrderRepository := &mocks.MockOrderRepository{}
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

func TestDeleteOrder_WhenOrderRepositoryDeleteMethodReturnsError_ReturnsError(t *testing.T) {
	//Given
	orderNumber := "1"
	mockOrderRepository := &mocks.MockOrderRepository{}
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

func getCreateOrderRequest() string {
	return `{
  "orderNumber": "1",
  "firstName": "Test",
  "lastName": "Sample",
  "totalAmount": 10.2,
  "address": "address",
  "city": "İstanbul",
  "district": "Bakırköy",
  "currencyCode": "TRY"
}`
}

func getUpdateOrderRequest() string {
	return `{
  "firstName": "Test",
  "lastName": "Sample",
  "totalAmount": 10.2,
  "address": "address",
  "city": "İstanbul",
  "district": "Bakırköy",
  "currencyCode": "TRY"
}`
}

var statusIdsThatNotBeValidForDeletion = []int{
	int(enum.Transferred),
	int(enum.Shipped),
	int(enum.Delivered),
}

var statusIdsThatNotBeValidForChangable = []int{
	int(enum.Transferred),
	int(enum.Shipped),
	int(enum.Delivered),
}
