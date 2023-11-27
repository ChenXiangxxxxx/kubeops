package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kube-backend/service"
	"net/http"
)

var Deployment deployment

type deployment struct{}

// 获取列表
func (d *deployment) GetDeployments(ctx *gin.Context) {
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
	data, err := service.Deployment.GetDeployments(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取deploy列表成功",
		"data": data,
	})
}

// 获取详情
func (d *deployment) GetDeploymentDetail(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		DeploymentName string `form:"deployment_name"`
		Namespace      string `form:"namespace"`
		Cluster        string `form:"cluster"`
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
	data, err := service.Deployment.GetDeploymentDetail(client, params.DeploymentName, params.Namespace)
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
		"msg":  "获取deploy详情成功",
		"data": data,
	})
}

// 删除
func (d *deployment) DeleteDeployment(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
		Cluster        string `json:"cluster"`
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
	err = service.Deployment.DeleteDeployment(client, params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除deploy详情成功",
		"data": nil,
	})
}

// 更新
func (d *deployment) UpdateDeployment(ctx *gin.Context) {
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
	err = service.Deployment.UpdateDeployment(client, params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新deploy详情成功",
		"data": nil,
	})
}

// 调整deploy副本数
func (d *deployment) ScaleDeployment(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
		ScaleNum       int    `json:"scale_num"`
		Cluster        string `json:"cluster"`
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
	data, err := service.Deployment.ScaleDeployment(client, params.DeploymentName, params.Namespace, params.ScaleNum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "调整deploy副本数成功",
		"data": data,
	})
}

// 重启deploy
func (d *deployment) RestartDeployment(ctx *gin.Context) {
	//接受参数,get请求为form格式。其他为json
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
		Cluster        string `json:"cluster"`
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
	err = service.Deployment.RestartDeployment(client, params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "重启deploy成功",
		"data": nil,
	})
}

// 创建deploy
func (d *deployment) CreateDeployment(ctx *gin.Context) {
	var (
		deployCreate = new(service.DeployCreate)
		err          error
	)
	//绑定参数
	//form格式使用ctx.Bind. json格式使用ctx.ShouldBindJson方法
	if err := ctx.ShouldBindJSON(deployCreate); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(deployCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	err = service.Deployment.CreateDeployment(client, deployCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建deploy成功",
		"data": nil,
	})
}
