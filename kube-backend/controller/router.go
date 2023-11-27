package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Router router

type router struct {
}

func (*router) InitApiRouter(r *gin.Engine) {
	r.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "testapi",
		})
	}).
		//pod
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		//deploy
		GET("/api/k8s/deployments", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		DELETE("/api/k8s/deployment/del", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).
		//service操作
		GET("/api/k8s/services", Servicev1.GetServices).
		POST("/api/k8s/service/create", Servicev1.CreateService).
		//ingress操作
		POST("/api/k8s/ingress/create", Ingress.CreateIngress).
		GET("/api/k8s/events", Event.GetList)

}
