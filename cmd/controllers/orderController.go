package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"simple-order-api/cmd/constants"
	"simple-order-api/cmd/helpers"
	"simple-order-api/cmd/models/request"
	"simple-order-api/cmd/models/response"
	"simple-order-api/cmd/services"
	"strings"
)

type OrderController struct {
	orderService services.OrderService
}

func NewOrderController(
	orderService services.OrderService,
) Controller {
	return &OrderController{
		orderService: orderService,
	}
}

// @Tags OrderController
// @Description Get Orders
// @Produce json
// @Success 200 {object} []response.Order
// @Failure 500 {object} response.ErrorResponse
// @Router /orders [get]
func (controller *OrderController) GetOrders() func(context *gin.Context) {
	return func(context *gin.Context) {
		orders, errorResp := controller.orderService.GetOrders()
		if errorResp != nil {
			context.JSON(http.StatusInternalServerError, errorResp)
			return
		}

		context.JSON(http.StatusOK, orders)
	}
}

// @Tags OrderController
// @Description Get Order By OrderNumber
// @Produce json
// @Success 200 {object} response.Order
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /orders/{orderNumber} [get]
// @Param orderNumber path string true "orderNumber"
func (controller *OrderController) GetOrderByOrderNumber() func(context *gin.Context) {
	return func(context *gin.Context) {
		orderNumber, orderNumberErr := getStringParam(context, constants.OrderNumber)
		if !helpers.IsValidString(orderNumber, orderNumberErr) {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.OrderNumberIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		order, errorResp := controller.orderService.GetOrder(orderNumber)
		if errorResp != nil {
			context.JSON(errorResp.StatusCode, errorResp)
			return
		}

		if order == nil {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusNotFound, constants.OrderNotFoundByOrderNumber).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		context.JSON(http.StatusOK, order)
	}
}

// @Tags OrderController
// @Description Create Order
// @Produce json
// @Success 201
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /orders [post]
// @Param request body request.CreateOrderRequest true "Create Order Request"
func (controller *OrderController) CreateOrder() func(context *gin.Context) {
	return func(context *gin.Context) {
		var createOrderRequest *request.CreateOrderRequest
		_ = mapstructure.Decode(getRequestBody(createOrderRequest, context), &createOrderRequest)

		if createOrderRequest == nil {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.CreateOrderRequestIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(createOrderRequest.OrderNumber)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.OrderNumberIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(createOrderRequest.FirstName)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.FirstNameIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(createOrderRequest.LastName)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.LastNameIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if createOrderRequest.TotalAmount <= 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.TotalAmountIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(createOrderRequest.Address)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.AddressIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(createOrderRequest.City)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.CityIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(createOrderRequest.District)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.DistrictIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(createOrderRequest.CurrencyCode)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.CurrencyCodeIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		createErr := controller.orderService.CreateOrder(*createOrderRequest)
		if createErr != nil {
			context.JSON(createErr.StatusCode, createErr)
			return
		}

		context.JSON(http.StatusCreated, "")
	}
}

// @Tags OrderController
// @Description Update Order
// @Produce json
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /orders/{orderNumber} [put]
// @Param orderNumber path string true "orderNumber"
// @Param request body request.UpdateOrderRequest true "Update Order Request"
func (controller *OrderController) UpdateOrder() func(context *gin.Context) {
	return func(context *gin.Context) {
		orderNumber, orderNumberErr := getStringParam(context, constants.OrderNumber)
		if !helpers.IsValidString(orderNumber, orderNumberErr) {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.OrderNumberIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		var updateOrderRequest *request.UpdateOrderRequest
		_ = mapstructure.Decode(getRequestBody(updateOrderRequest, context), &updateOrderRequest)

		if updateOrderRequest == nil {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.UpdateOrderRequestIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(updateOrderRequest.FirstName)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.FirstNameIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(updateOrderRequest.LastName)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.LastNameIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if updateOrderRequest.TotalAmount <= 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.TotalAmountIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(updateOrderRequest.Address)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.AddressIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(updateOrderRequest.City)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.CityIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(updateOrderRequest.District)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.DistrictIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		if len(strings.TrimSpace(updateOrderRequest.CurrencyCode)) == 0 {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.CurrencyCodeIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		createErr := controller.orderService.UpdateOrder(orderNumber, *updateOrderRequest)
		if createErr != nil {
			context.JSON(createErr.StatusCode, createErr)
			return
		}

		context.JSON(http.StatusNoContent, "")
	}
}

// @Tags OrderController
// @Description Delete Order
// @Produce json
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /orders/{orderNumber} [delete]
// @Param orderNumber path string true "orderNumber"
func (controller *OrderController) DeleteOrder() func(context *gin.Context) {
	return func(context *gin.Context) {
		orderNumber, orderNumberErr := getStringParam(context, constants.OrderNumber)
		if !helpers.IsValidString(orderNumber, orderNumberErr) {
			errorResponse := response.NewErrorBuilder().
				SetError(http.StatusBadRequest, constants.OrderNumberIsNotValid).
				Build()
			context.JSON(errorResponse.StatusCode, errorResponse)
			return
		}

		deleteErr := controller.orderService.DeleteOrder(orderNumber)
		if deleteErr != nil {
			context.JSON(deleteErr.StatusCode, deleteErr)
			return
		}

		context.JSON(http.StatusNoContent, "")
	}
}

func (controller *OrderController) Register(engine *gin.Engine) {
	engine.GET("/orders", controller.GetOrders())
	engine.GET("/orders/:orderNumber", controller.GetOrderByOrderNumber())
	engine.POST("/orders", controller.CreateOrder())
	engine.PUT("/orders/:orderNumber", controller.UpdateOrder())
	engine.DELETE("/orders/:orderNumber", controller.DeleteOrder())
}
