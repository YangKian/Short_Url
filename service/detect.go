package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"shortUrl/models"
)

func HealthCheck(c *gin.Context) {
	if clickCounter.flag == true {
		urlCode := models.UrlCode{}
		urlCode.SaveClicks(clickCounter.counterMap)
		log.Println("save the clickCounter to db.")
		clickCounter.flag = false
	}
	message := "OK"
	c.String(http.StatusOK, "\n"+message)
}
