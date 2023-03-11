package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type SwaggerControllerImp struct {
}

func NewSwaggerController() Controller {
	return &SwaggerControllerImp{}
}

// @Param Authorization header string true "With the bearer started"
func (controller *SwaggerControllerImp) Register(engine *gin.Engine) {
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
