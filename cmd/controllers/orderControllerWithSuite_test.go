package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"simple-order-api/cmd/constants"
	"simple-order-api/cmd/mocks"
	"simple-order-api/cmd/models/request"
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

func (o *OrderControllerSuite) readError() response.ErrorResponse {
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(o.recorder.Body.Bytes(), &errResponse)
	return errResponse
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
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()
	o.mockOrderService.On("GetOrders").Return(nil, &serviceErr)

	//When
	o.sendRequest("GET", "/orders", nil)

	//Then
	assert.Equal(o.T(), http.StatusInternalServerError, o.recorder.Code)
	assert.Equal(o.T(), serviceErr, o.readError())
}

func (o *OrderControllerSuite) TestGetOrderByOrderNumberWithSuite() {
	//Given
	order := response.Order{

		OrderNumber:  "1",
		FirstName:    "Ahmet",
		LastName:     "Ata",
		TotalAmount:  121.13,
		Address:      "Lorem ipsum dolor sit amet",
		City:         "İstanbul",
		District:     "Silivri",
		StatusId:     2,
		CurrencyCode: "TR",
	}
	o.mockOrderService.On("GetOrder", mock.Anything).Return(&order, nil)

	//When
	o.sendRequest("GET", "/orders/123456", nil)

	//Then
	assert.Equal(o.T(), http.StatusOK, o.recorder.Code)
	assert.NotNil(o.T(), o.recorder.Body)
	expectedResp := &response.Order{}
	_ = json.Unmarshal(o.recorder.Body.Bytes(), expectedResp)
	assert.Equal(o.T(), order, *expectedResp)
	o.mockOrderService.AssertCalled(o.T(), "GetOrder", "123456")
	o.mockOrderService.AssertNumberOfCalls(o.T(), "GetOrder", 1)
}

func (o *OrderControllerSuite) TestGetOrderByOrderNumberWithSuite_WhenOrderNumberIsInvalid_returnsBadRequestError() {
	//Given
	//When
	o.sendRequest("GET", "/orders/%20", nil)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.OrderNumberIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "GetOrder", 0)

}

func (o *OrderControllerSuite) TestGetOrderByOrderNumberWithSuite_WhenOrderNotFound_returnsNotFoundError() {
	//Given
	o.mockOrderService.On("GetOrder", mock.Anything).Return(nil, nil)

	//When
	o.sendRequest("GET", "/orders/123456", nil)

	//Then
	assert.Equal(o.T(), http.StatusNotFound, o.recorder.Code)
	assert.Equal(o.T(), constants.OrderNotFoundByOrderNumber, o.readError().Message)
}

func (o *OrderControllerSuite) TestGetOrderByOrderNumberWithSuite_WhenServiceReturnsError_ReturnsInternalServerError() {
	//Given
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()
	o.mockOrderService.On("GetOrder", mock.Anything).Return(nil, &serviceErr)

	//When
	o.sendRequest("GET", "/orders/123456", nil)

	//Then
	assert.Equal(o.T(), http.StatusInternalServerError, o.recorder.Code)
	assert.Equal(o.T(), serviceErr, o.readError())
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite() {
	//Given
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusCreated, o.recorder.Code)
	var createOrderRequest = request.CreateOrderRequest{
		OrderNumber:  "1",
		FirstName:    "Test",
		LastName:     "Sample",
		TotalAmount:  10.2,
		Address:      "address",
		City:         "İstanbul",
		District:     "Bakırköy",
		CurrencyCode: "TRY",
	}
	o.mockOrderService.AssertCalled(o.T(), "CreateOrder", createOrderRequest)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 1)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenRequestIsNotValid_ReturnsBadRequest() {
	//Given
	//When
	o.sendRequest("POST", "/orders", nil)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.CreateOrderRequestIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 0)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenOrderNumberIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.OrderNumber = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.OrderNumberIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 0)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenFirstNameIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.FirstName = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.FirstNameIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 0)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenLastNameIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.LastName = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.LastNameIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 0)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenTotalAmountIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.TotalAmount = -12.13
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.TotalAmountIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 0)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenAddressIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.Address = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.AddressIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 0)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenCityIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.City = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.CityIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 0)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenDistrictIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.District = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.DistrictIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 0)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenCurrencyCodeIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.CurrencyCode = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.CurrencyCodeIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "CreateOrder", 0)
}

func (o *OrderControllerSuite) TestCreateOrderWithSuite_WhenOrderServiceReturnsError_returnsError() {
	//Given
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusConflict, "notFound").
		Build()
	o.mockOrderService.On("CreateOrder", mock.Anything).Return(&serviceErr)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("POST", "/orders", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusConflict, o.recorder.Code)
	assert.Equal(o.T(), serviceErr, o.readError())
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite() {
	//Given
	o.mockOrderService.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/123456", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusNoContent, o.recorder.Code)
	var updateOrderRequest = request.UpdateOrderRequest{
		FirstName:    "Test",
		LastName:     "Sample",
		TotalAmount:  10.2,
		Address:      "address",
		City:         "İstanbul",
		District:     "Bakırköy",
		CurrencyCode: "TRY",
	}
	o.mockOrderService.AssertCalled(o.T(), "UpdateOrder", "123456", updateOrderRequest)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 1)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenRequestIsNotValid_ReturnsBadRequest() {
	//Given
	//When
	o.sendRequest("PUT", "/orders/123456", nil)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.UpdateOrderRequestIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 0)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenOrderNumberIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("UpdateOrder", mock.Anything).Return(nil)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/%20", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.OrderNumberIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 0)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenFirstNameIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("UpdateOrder", mock.Anything).Return(nil)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	serviceReq.FirstName = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/123456", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.FirstNameIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 0)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenLastNameIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("UpdateOrder", mock.Anything).Return(nil)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	serviceReq.LastName = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/123456", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.LastNameIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 0)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenTotalAmountIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("UpdateOrder", mock.Anything).Return(nil)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	serviceReq.TotalAmount = -12.13
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/123456", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.TotalAmountIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 0)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenAddressIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("UpdateOrder", mock.Anything).Return(nil)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	serviceReq.Address = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/123456", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.AddressIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 0)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenCityIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("UpdateOrder", mock.Anything).Return(nil)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	serviceReq.City = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/123456", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.CityIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 0)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenDistrictIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("UpdateOrder", mock.Anything).Return(nil)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	serviceReq.District = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/123456", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.DistrictIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 0)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenCurrencyCodeIsNotValid_ReturnsBadRequest() {
	//Given
	o.mockOrderService.On("UpdateOrder", mock.Anything).Return(nil)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	serviceReq.CurrencyCode = ""
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/123456", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.CurrencyCodeIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "UpdateOrder", 0)
}

func (o *OrderControllerSuite) TestUpdateOrderWithSuite_WhenOrderServiceReturnsError_returnsError() {
	//Given
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusNotFound, "notFound").
		Build()
	o.mockOrderService.On("UpdateOrder", mock.Anything, mock.Anything).Return(&serviceErr)
	serviceReq := &request.UpdateOrderRequest{}
	_ = json.Unmarshal([]byte(getUpdateOrderRequest()), serviceReq)
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	o.sendRequest("PUT", "/orders/123456", reqBodyBytes)

	//Then
	assert.Equal(o.T(), http.StatusNotFound, o.recorder.Code)
	assert.Equal(o.T(), serviceErr, o.readError())
}

func (o *OrderControllerSuite) TestDeleteOrderWithSuite() {
	//Given
	o.mockOrderService.On("DeleteOrder", mock.Anything).Return(nil)

	//When
	o.sendRequest("DELETE", "/orders/123456", nil)

	//Then
	assert.Equal(o.T(), http.StatusNoContent, o.recorder.Code)
	o.mockOrderService.AssertCalled(o.T(), "DeleteOrder", "123456")
	o.mockOrderService.AssertNumberOfCalls(o.T(), "DeleteOrder", 1)
}

func (o *OrderControllerSuite) TestDeleteOrderWithSuite_WhenOrderNumberIsInvalid_returnsBadRequestError() {
	//Given
	//When
	o.sendRequest("DELETE", "/orders/%20", nil)

	//Then
	assert.Equal(o.T(), http.StatusBadRequest, o.recorder.Code)
	assert.Equal(o.T(), constants.OrderNumberIsNotValid, o.readError().Message)
	o.mockOrderService.AssertNumberOfCalls(o.T(), "DeleteOrder", 0)
}

func (o *OrderControllerSuite) TestDeleteOrderWithSuite_WhenOrderServiceReturnsError_returnsError() {
	//Given
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusNotFound, "notFound").
		Build()
	o.mockOrderService.On("DeleteOrder", mock.Anything).Return(&serviceErr)

	//When
	o.sendRequest("DELETE", "/orders/123456", nil)

	//Then
	assert.Equal(o.T(), http.StatusNotFound, o.recorder.Code)
	assert.Equal(o.T(), serviceErr, o.readError())
}
