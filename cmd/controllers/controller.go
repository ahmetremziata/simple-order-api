package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type Controller interface {
	Register(engine *gin.Engine)
}

func getStringParam(context *gin.Context, paramName string) (string, error) {
	valueParam := context.Param(paramName)
	return valueParam, nil
}

func getRequestBody(request interface{}, context *gin.Context) interface{} {
	if context.Request.Body == nil {
		return nil
	}

	byteBody, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		return nil
	}

	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(byteBody))
	err = context.BindJSON(&request)

	if err != nil {
		return nil
	}

	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(byteBody))
	return request
}
