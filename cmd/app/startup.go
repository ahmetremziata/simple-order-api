package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	controllers2 "simple-order-api/cmd/controllers"
	"simple-order-api/cmd/docs"
	"simple-order-api/cmd/models"
	"simple-order-api/cmd/services"
)

func StartServer() {
	serverConfig := models.ServerConfig{
		Port: ":8080",
		Host: "localhost:8080",
	}
	docs.SwaggerInfo.Host = serverConfig.Host
	engine := setHttpServerConfigs()
	orderService := services.NewOrderService()
	orderController := controllers2.NewOrderController(orderService)
	swaggerController := controllers2.NewSwaggerController()
	swaggerController.Register(engine)
	orderController.Register(engine)

	if err := engine.Run(serverConfig.Port); err != nil {
		fmt.Println("An error has occured while starting web server")
		panic(true)
	}

	fmt.Println("Order web server started successfully!")

}

func setHttpServerConfigs() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(corsMiddleware())
	return engine
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
