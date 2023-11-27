package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct{}

type DeploymentResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

// 定义deploy创建的结构体
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	Cluster       string            `json:"cluster"`
}

// 获取dep列表
func (d *deployment) GetDeployments(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (deploymentResp *DeploymentResp, err error) {
	deploymentList, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取deploy列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取deploy列表失败,%v\n", err))
	}
	selectableData := &dataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
		dataSelectorQuery: &DataSelectorQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤，再排序和分页
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()
	deployments := d.fromCells(data.GenericDataList)
	return &DeploymentResp{
		Items: deployments,
		Total: total,
	}, nil
}

// 删除deployment
func (d *deployment) DeleteDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	err = client.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("删除Deployment失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除Deployment失败, %v\n", err))
	}

	return nil
}

// 更新deployment
func (d *deployment) UpdateDeployment(client *kubernetes.Clientset, namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}

	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}

	_, err = client.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("更新Deployment失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新Deployment失败, %v\n", err))
	}
	return nil
}

// 获取deployment详情
func (d *deployment) GetDeploymentDetail(client *kubernetes.Clientset, deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取Deployment详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取Deployment详情失败, %v\n", err))
	}

	return deployment, nil
}

// 修改deploy副本数
func (d *deployment) ScaleDeployment(client *kubernetes.Clientset, deploymentName, namespace string, scaleNum int) (replica int32, err error) {
	scale, err := client.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取Deployment的副本信息失败, %v\n", err))
		return 0, errors.New(fmt.Sprintf("获取Deployment的副本信息失败, %v\n", err))
	}
	//修改副本数
	scale.Spec.Replicas = int32(scaleNum)
	//更新副本数，传入scale
	newScale, err := client.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("更新Deployment详情失败, %v\n", err))
		return 0, errors.New(fmt.Sprintf("更新Deployment详情失败, %v\n", err))
	}
	return newScale.Spec.Replicas, nil
}

// 重启deploy
func (d *deployment) RestartDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	//通过patch方法实现重启
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{"name": deploymentName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}
	//序列化成字符串
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("序列化失败, %v\n", err))
	}
	//调用patch方法更新
	_, err = client.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName,
		"application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("重启Deployment失败, %v\n", err))
		return errors.New(fmt.Sprintf("重启Deployment失败, %v\n", err))
	}
	return nil
}

// 创建deploy
func (d *deployment) CreateDeployment(client *kubernetes.Clientset, data *DeployCreate) (err error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}
	//判断是否打开监控检查功能，若打开则定义ReadinessProbe和LivenessProbe
	if data.HealthCheck {
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					//intstr.IntOrString的作用是端口可以定义为整行，也可以定义为字符串
					//type=0表示整行，使用intVal
					//type=1表示字符串，使用strVal
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			//初始化等待时间
			InitialDelaySeconds: 15,
			//超时时间
			TimeoutSeconds: 15,
			//执行间隔
			PeriodSeconds: 5,
		}
	}
	//定义容器的limit和request资源
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits =
		map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests =
		map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}
	//创建deployment
	_, err = client.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("创建Deployment失败, %v\n", err))
		return errors.New(fmt.Sprintf("创建Deployment失败, %v\n", err))
	}
	return nil
}

// 定义DataCell到deployment类型转换的方法
func (p *deployment) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i]) //这里是std从podCell的结构体转位DataCell的接口类型。
	}
	return cells
}

func (p *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell)) //cells是接口，断言成podCell类型，再通过corev1.Pod转换。
	}
	return deployments
}
