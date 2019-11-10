package main

import (
	"MyProject/Short_Url/models"
	"MyProject/Short_Url/pkg/setting"
	"MyProject/Short_Url/router"
	"fmt"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	// log.SetFlags(log.Ldate | log.Lshortfile)
	setting.Start()
	models.Start()

}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := router.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := &http.Server{
		Addr:         endPoint,
		Handler:      routersInit,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
}
