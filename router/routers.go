package routers

import (
	"MyProject/Short_Url/pkg/service"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	urlService := service.NewShortUrlService()
	api := r.Group("/api/v1")
	{
		api.POST("/create", urlService.Create)
		// api.POST("/get", Shorturl.Get)
	}
}
