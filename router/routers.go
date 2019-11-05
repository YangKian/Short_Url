package router

import (
	"MyProject/Short_Url/pkg/service"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	urlService := service.NewShortUrlService()
	api := r.Group("/api/v1")
	{
		api.POST("/create", urlService.Create)
		// api.POST("/get", Shorturl.Get)
	}
	return r
}
