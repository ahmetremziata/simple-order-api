package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"simple-order-api/cmd/constants"
	"simple-order-api/cmd/mocks"
	"simple-order-api/cmd/models/request"
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
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()

	mockOrderService.On("GetOrders").Return(nil, &serviceErr)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("GET", "/orders", nil)
	engine.ServeHTTP(w, req)

	//Then
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, errResponse, serviceErr)
}

func TestGetOrderByOrderNumber(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
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
	mockOrderService.On("GetOrder", mock.Anything).Return(&order, nil)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("GET", "/orders/123456", nil)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)
	expectedResp := &response.Order{}
	_ = json.Unmarshal(w.Body.Bytes(), expectedResp)
	assert.Equal(t, order, *expectedResp)
	mockOrderService.AssertCalled(t, "GetOrder", "123456")
	mockOrderService.AssertNumberOfCalls(t, "GetOrder", 1)
}

func TestGetOrderByOrderNumber_WhenOrderNumberIsInvalid_returnsBadRequestError(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("GET", "/orders/%20", nil)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.OrderNumberIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "GetOrder", 0)
}

func TestGetOrderByOrderNumber_WhenOrderNotFound_returnsNotFoundError(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("GetOrder", mock.Anything).Return(nil, nil)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("GET", "/orders/123456", nil)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusNotFound, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.OrderNotFoundByOrderNumber, errResponse.Message)
}

func TestGetOrderByOrderNumber_WhenOrderServiceReturnsError_returnsInternalServerError(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()
	mockOrderService.On("GetOrder", mock.Anything).Return(nil, &serviceErr)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("GET", "/orders/123456", nil)
	engine.ServeHTTP(w, req)

	//Then
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, errResponse, serviceErr)
}

func TestCreateOrder(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusCreated, w.Code)
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
	mockOrderService.AssertCalled(t, "CreateOrder", createOrderRequest)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 1)
}

func TestCreateOrder_WhenRequestIsNotValid_ReturnsBadRequest(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("POST", "/orders", nil)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.CreateOrderRequestIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 0)
}

func TestCreateOrder_WhenOrderNumberIsNotValid_ReturnsBadRequest(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.OrderNumber = ""
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.OrderNumberIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 0)
}

func TestCreateOrder_WhenFirstNameIsNotValid_ReturnsBadRequest(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.FirstName = ""
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.FirstNameIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 0)
}

func TestCreateOrder_WhenLastNameIsNotValid_ReturnsBadRequest(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.LastName = ""
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.LastNameIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 0)
}

func TestCreateOrder_WhenTotalAmountIsNotValid_ReturnsBadRequest(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.TotalAmount = 0
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.TotalAmountIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 0)
}

func TestCreateOrder_WhenAddressIsNotValid_ReturnsBadRequest(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.Address = ""
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.AddressIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 0)
}

func TestCreateOrder_WhenCityIsNotValid_ReturnsBadRequest(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.City = ""
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.CityIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 0)
}

func TestCreateOrder_WhenDistrictIsNotValid_ReturnsBadRequest(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.District = ""
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.DistrictIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 0)
}

func TestCreateOrder_WhenCurrencyCodeIsNotValid_ReturnsBadRequest(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("CreateOrder", mock.Anything).Return(nil)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	serviceReq.CurrencyCode = ""
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.CurrencyCodeIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "CreateOrder", 0)
}

func TestCreateOrder_WhenOrderServiceReturnsError_returnsInternalServerError(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()
	mockOrderService.On("CreateOrder", mock.Anything).Return(&serviceErr)
	serviceReq := &request.CreateOrderRequest{}
	_ = json.Unmarshal([]byte(getCreateOrderRequest()), serviceReq)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()
	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(*serviceReq)

	//When
	req, _ := http.NewRequest("POST", "/orders", reqBodyBytes)
	engine.ServeHTTP(w, req)

	//Then
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, errResponse, serviceErr)
}

func TestDeleteOrder(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	mockOrderService.On("DeleteOrder", mock.Anything).Return(nil)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("DELETE", "/orders/123456", nil)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockOrderService.AssertCalled(t, "DeleteOrder", "123456")
	mockOrderService.AssertNumberOfCalls(t, "DeleteOrder", 1)
}

func TestDeleteOrder_WhenOrderNumberIsInvalid_returnsBadRequestError(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("DELETE", "/orders/%20", nil)
	engine.ServeHTTP(w, req)

	//Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, constants.OrderNumberIsNotValid, errResponse.Message)
	mockOrderService.AssertNumberOfCalls(t, "DeleteOrder", 0)
}

func TestDeleteOrder_WhenOrderServiceReturnsError_returnsInternalServerError(t *testing.T) {
	//Given
	engine := gin.New()
	mockOrderService := &mocks.FakeOrderService{}
	serviceErr := response.NewErrorBuilder().
		SetError(http.StatusInternalServerError, "test").
		Build()
	mockOrderService.On("DeleteOrder", mock.Anything).Return(&serviceErr)
	controller := NewOrderController(mockOrderService)
	controller.Register(engine)
	w := httptest.NewRecorder()

	//When
	req, _ := http.NewRequest("DELETE", "/orders/123456", nil)
	engine.ServeHTTP(w, req)

	//Then
	errResponse := response.ErrorResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, errResponse, serviceErr)
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
