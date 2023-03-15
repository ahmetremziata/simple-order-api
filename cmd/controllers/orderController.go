package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-order-api/cmd/constants"
	"simple-order-api/cmd/helpers"
	"simple-order-api/cmd/models/response"
	"simple-order-api/cmd/services"
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
		orderNumber, orderNumberErr := GetStringParam(context, constants.OrderNumber)
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
		orderNumber, orderNumberErr := GetStringParam(context, constants.OrderNumber)
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
	engine.DELETE("/orders/:orderNumber", controller.DeleteOrder())
}
