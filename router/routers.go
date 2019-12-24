package router

import (
	"github.com/gin-gonic/gin"
	"shortUrl/service"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	urlService := service.NewShortUrlService()
	api := r.Group("/")
	{
		api.POST("/url", urlService.SingleCreate) //传入url生成shortUrl
		api.POST("/shortUrl", urlService.TransToUrl) //传入shortUrl返回url
	}

	detecting := r.Group("/detect")
	{
		detecting.GET("/health", service.HealthCheck)
	}
	return r
}
