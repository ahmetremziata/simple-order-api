package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-order-api/cmd/constants"
	"simple-order-api/cmd/helpers"
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
// @Failure 500 {object} error
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
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} error
// @Router /orders/{orderNumber} [get]
// @Param orderNumber path string true "orderNumber"
func (controller *OrderController) GetOrderByOrderNumber() func(context *gin.Context) {
	return func(context *gin.Context) {
		orderNumber, orderNumberErr := GetStringParam(context, constants.OrderNumber)
		if !helpers.IsValidString(orderNumber, orderNumberErr) {
			context.JSON(http.StatusBadRequest, constants.OrderNumberIsNotValid)
			return
		}

		order, errorResp := controller.orderService.GetOrder(orderNumber)
		if errorResp != nil {
			context.JSON(http.StatusInternalServerError, errorResp)
			return
		}

		if order == nil {
			context.JSON(http.StatusNotFound, constants.OrderNotFoundByOrderNumber)
			return
		}

		context.JSON(http.StatusOK, order)
	}
}

func (controller *OrderController) Register(engine *gin.Engine) {
	engine.GET("/orders", controller.GetOrders())
	engine.GET("/orders/:orderNumber", controller.GetOrderByOrderNumber())
}
