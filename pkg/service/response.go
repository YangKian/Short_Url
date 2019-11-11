package service

import (
	"MyProject/Short_Url/contants"
	"net/http"

	"github.com/gin-gonic/gin"
)

//type Response struct{}

func success(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  contants.MsgGeter(code),
		"data": data,
	})
}

func innerFail(c *gin.Context, code int) {
	c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
		"code": code,
		"msg":  contants.MsgGeter(code),
		"data": "",
	})
}

func requestFail(c *gin.Context, code int) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"code": code,
		"msg":  contants.MsgGeter(code),
		"data": "",
	})
}

func redirect(c *gin.Context, path string) {
	c.Redirect(http.StatusFound, path)
}
