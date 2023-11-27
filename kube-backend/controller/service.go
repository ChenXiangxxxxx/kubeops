package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kube-backend/service"
	"net/http"
)

var Servicev1 servicev1

type servicev1 struct{}

func (s servicev1) GetServices(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Cluster    string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
		return
	}
	client, _ := service.K8s.GetClient(params.Cluster)
	data, err := service.Servicev1.GetServices(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取service列表成功",
		"data": data,
	})
}

// 创建deployment
func (s *servicev1) CreateService(ctx *gin.Context) {
	var (
		serviceCreate = new(service.ServiceCreate)
		err           error
	)
	//绑定参数
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.ShouldBindJSON(serviceCreate); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败, %v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败, %v\n", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(serviceCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	err = service.Servicev1.CreateService(client, serviceCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			//"code": 90500, //业务状态
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建Service成功",
		"data": nil,
	})
}
