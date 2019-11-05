package service
import (
	"MyProject/Short_Url/contants"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct{}

func (r *Response) success(c *gin.Context, data map[string]interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": contants.SUCCESS,
		"msg":  "ok",
		"data": data,
	})
}

func (r *Response) fail(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": "",
	})
}
