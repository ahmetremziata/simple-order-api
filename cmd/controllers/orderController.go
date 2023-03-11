package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

func (controller *OrderController) Register(engine *gin.Engine) {
	engine.GET("/orders", controller.GetOrders())
}
