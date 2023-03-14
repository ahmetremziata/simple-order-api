package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"simple-order-api/cmd/mocks"
	"simple-order-api/cmd/models/response"
	"testing"
)

type OrderControllerSuite struct {
	suite.Suite
	orderController  Controller
	engine           *gin.Engine
	mockOrderService *mocks.FakeOrderService
	recorder         *httptest.ResponseRecorder
}

func TestOrderControllerInit(t *testing.T) {
	suite.Run(t, new(OrderControllerSuite))
}

func (o *OrderControllerSuite) SetupTest() {
	o.engine = gin.New()
	o.mockOrderService = new(mocks.FakeOrderService)
	o.orderController = NewOrderController(o.mockOrderService)
	o.orderController.Register(o.engine)
	o.recorder = httptest.NewRecorder()
}

func (o *OrderControllerSuite) sendRequest(method, uri string, body io.Reader) {
	req, httpEr := http.NewRequest(method, uri, body)
	require.NoError(o.T(), httpEr)
	o.engine.ServeHTTP(o.recorder, req)
}

func (o *OrderControllerSuite) TestGetOrdersWithSuite() {
	//Given
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
	o.mockOrderService.On("GetOrders").Return(orders, nil)

	//When
	o.sendRequest("GET", "/orders", nil)

	//Then
	assert.Equal(o.T(), http.StatusOK, o.recorder.Code)
	assert.NotNil(o.T(), o.recorder.Body)
	expectedResp := &[]response.Order{}
	_ = json.Unmarshal(o.recorder.Body.Bytes(), expectedResp)
	assert.Equal(o.T(), orders, *expectedResp)
	o.mockOrderService.AssertCalled(o.T(), "GetOrders")
	o.mockOrderService.AssertNumberOfCalls(o.T(), "GetOrders", 1)
}

func (o *OrderControllerSuite) TestGetOrdersWithSuite_WhenServiceReturnsError_ReturnsInternalServerError() {
	//Given
	serviceErr := errors.New("service Error")
	o.mockOrderService.On("GetOrders").Return(nil, &serviceErr)

	//When
	o.sendRequest("GET", "/orders", nil)

	//Then
	assert.Equal(o.T(), http.StatusInternalServerError, o.recorder.Code)
}
