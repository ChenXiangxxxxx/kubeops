package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kube-backend/config"
	"kube-backend/controller"
	"kube-backend/db"
	"kube-backend/middle"
	"kube-backend/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(middle.Cors())
	//r.Use(middle.JWTAuth())
	//初始化数据库
	db.Init()
	service.K8s.Init()
	controller.Router.InitApiRouter(r)
	//启动task
	go func() {
		service.Event.WatchEventTask("TST-1")
	}()

	//websocket 启动
	wsHandler := http.NewServeMux()
	wsHandler.HandleFunc("/ws", service.Terminal.WsHandler)
	ws := &http.Server{
		Addr:    config.WsAddr,
		Handler: wsHandler,
	}
	go func() {
		if err := ws.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}
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

	//关闭websocket
	if err := ws.Shutdown(ctx); err != nil {
		log.Fatal("websocket 关闭异常：", err)
	}
	log.Println("websocket退出成功")

	//关闭ginserver
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("gin server关闭异常", err)
	}
	logger.Info("gin server退出成功")

	if err := db.Close(); err != nil {
		log.Fatal("DB关闭异常", err)
	}

}
