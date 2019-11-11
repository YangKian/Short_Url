package router

import (
	"MyProject/Short_Url/pkg/service"
	//ginSwagger "github.com/swaggo/gin-swagger"
	//"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	urlService := service.NewShortUrlService()
	api := r.Group("/api/v1")
	{
		api.POST("/create", urlService.SingleCreate)
		api.POST("/restore", urlService.TransToUrl)
		// api.POST("/get", Shorturl.Get)
	}
	return r
}
