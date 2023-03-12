package controllers

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Register(engine *gin.Engine)
}

func GetStringParam(context *gin.Context, paramName string) (string, error) {
	valueParam := context.Param(paramName)
	return valueParam, nil
}
