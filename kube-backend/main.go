package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kube-backend/config"
	"kube-backend/controller"
	"kube-backend/db"
	"kube-backend/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	r := gin.Default()
	//初始化数据库
	db.Init()
	service.K8s.Init()
	controller.Router.InitApiRouter(r)
	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}
	//启动task
	go func() {
		service.Event.WatchEventTask("TST-1")
	}()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen :%s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("gin server关闭异常", err)
	}
	logger.Info("gin server退出成功")

	if err := db.Close(); err != nil {
		log.Fatal("DB关闭异常", err)
	}

}
