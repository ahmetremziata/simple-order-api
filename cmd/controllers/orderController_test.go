package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"simple-order-api/cmd/models/response"
	"testing"
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

func TestGetOrders(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &FakeOrderService{}
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
	mockOrderService.On("GetOrders").Return(orders, nil)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("GET", "/orders", nil)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)
	expectedResp := &[]response.Order{}
	_ = json.Unmarshal(w.Body.Bytes(), expectedResp)
	assert.Equal(t, orders, *expectedResp)
	mockOrderService.AssertCalled(t, "GetOrders")
	mockOrderService.AssertNumberOfCalls(t, "GetOrders", 1)
}

func TestGetOrders_WhenServiceReturnsError_ReturnsInternalServerError(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &FakeOrderService{}
	serviceErr := errors.New("service Error")

	mockOrderService.On("GetOrders").Return(nil, &serviceErr)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("GET", "/orders", nil)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
