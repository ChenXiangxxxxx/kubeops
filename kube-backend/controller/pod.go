package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kube-backend/service"
	"net/http"
)

var Pod pod

type pod struct{}

// 获取pod列表
func (p *pod) GetPods(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Cluster    string `form:"cluster"`
	})
	//绑定参数
	//form格式使用ctx.Bind. json格式使用ctx.ShouldBindJson方法
	if err := ctx.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	data, err := service.Pod.GetPods(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
	})
}

// 获取pod详情
func (p *pod) GetPodDetail(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
	})
	//绑定参数
	//form格式使用ctx.Bind. json格式使用ctx.ShouldBindJson方法
	if err := ctx.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	data, err := service.Pod.GetPodDetail(client, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//byte, _ := json.Marshal(data)
	//ctx.JSON(http.StatusOK, gin.H{
	//	"msg":  "获取pod详情成功",
	//	"data": string(byte),
	//})
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod详情成功",
		"data": data,
	})
}

// 删除pod
func (p *pod) DeletePod(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
	})
	//绑定参数
	//form格式使用ctx.Bind. json格式使用ctx.ShouldBindJson方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	err = service.Pod.DeletePod(client, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除pod详情成功",
		"data": nil,
	})
}

// 更新pod
func (p *pod) UpdatePod(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
		Content   string `json:"content"`
	})
	//绑定参数
	//form格式使用ctx.Bind. json格式使用ctx.ShouldBindJson方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	err = service.Pod.UpdatePod(client, params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新pod详情成功",
		"data": nil,
	})
}

// 获取容器
func (p *pod) GetPodContainer(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
	})
	//绑定参数
	//form格式使用ctx.Bind. json格式使用ctx.ShouldBindJson方法
	if err := ctx.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	data, err := service.Pod.GetPodContainer(client, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器成功",
		"data": data,
	})
}

// 获取容器日志
func (p *pod) GetPodLog(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		ContainerName string `form:"container_name"`
		PodName       string `form:"pod_name"`
		Namespace     string `form:"namespace"`
		Cluster       string `form:"cluster"`
	})
	//绑定参数
	//form格式使用ctx.Bind. json格式使用ctx.ShouldBindJson方法
	if err := ctx.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	data, err := service.Pod.GetPodLog(client, params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器日志成功",
		"data": data,
	})
}
