package routers

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api/v1")
	{
		api.POST("/create", Shorturl.Create)
		api.POST("/get", Shorturl.Get)
	}
}
