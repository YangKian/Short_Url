package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"shortUrl/models"
	"shortUrl/pkg/setting"
	"shortUrl/router"
	"shortUrl/service"
	"time"
)

func init() {
	setting.Start()
	models.Start()
	service.Start()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := router.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeoutTimer
	writeTimeout := setting.ServerSetting.WriteTimeoutTimer
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := &http.Server{
		Addr:         endPoint,
		Handler:      routersInit,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	go func() {
		for {
			if err := pingServer(server); err != nil {
				log.Fatalf("heart beat detect has no response, err: %v\n", err)
			}
			time.Sleep(setting.ServerSetting.HeartBeatCheckTimer)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("[Main]: start service failed, err: %v\n", err)
	}
}

func pingServer(svr *http.Server) error {
	for i := 0; i < 10; i++ {
		resp, err := http.Get("http://127.0.0.1" + svr.Addr + "/detect/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		log.Println("[PingServer]: Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("[PingServer]: cannot connect to the router")
}
