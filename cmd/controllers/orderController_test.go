package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"simple-order-api/cmd/mocks"
	"simple-order-api/cmd/models/response"
	"testing"
)

func TestGetOrders(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
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
	mockOrderService := &mocks.FakeOrderService{}
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
