package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Ingress ingress

type ingress struct{}

type IngressesResp struct {
	Items []nwv1.Ingress `json:"items"`
	Total int            `json:"total"`
}

// 定义ServiceCreate结构体，用于创建service需要的参数属性的定义
type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Label     map[string]string      `json:"label"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
	Cluster   string                 `json:"cluster"`
}

// 定义ingress的path结构体
type HttpPath struct {
	Path        string        `json:"path"`
	PathType    nwv1.PathType `json:"path_type"`
	ServiceName string        `json:"service_name"`
	ServicePort int32         `json:"service_port"`
}

// 获取ingress列表，支持过滤、排序、分页
func (i *ingress) GetIngresses(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (ingressesResp *IngressesResp, err error) {
	//获取ingressList类型的ingress列表
	ingressList, err := client.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Ingress列表失败, " + err.Error()))
		return nil, errors.New("获取Ingress列表失败, " + err.Error())
	}
	//将ingressList中的ingress列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: i.toCells(ingressList.Items),
		dataSelectorQuery: &DataSelectorQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()

	//将[]DataCell类型的ingress列表转为v1.ingress列表
	ingresss := i.fromCells(data.GenericDataList)

	return &IngressesResp{
		Items: ingresss,
		Total: total,
	}, nil
}

// 获取ingress详情
func (i *ingress) GetIngresstDetail(client *kubernetes.Clientset, ingressName, namespace string) (ingress *nwv1.Ingress, err error) {
	ingress, err = client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Ingress详情失败, " + err.Error()))
		return nil, errors.New("获取Ingress详情失败, " + err.Error())
	}

	return ingress, nil
}
func (i *ingress) CreateIngress(client *kubernetes.Clientset, data *IngressCreate) (err error) {
	//声明nwv1.IngressRule和nwv1.HTTPIngressPath变量，后面用于数据组装
	//ingressRule代表的是Hosts
	var ingressRules = make([]nwv1.IngressRule, 0)
	//httpIngressPaths代表的是Paths
	var httpIngressPaths = make([]nwv1.HTTPIngressPath, 0)
	//将data中的数据组装成nwv1.Ingress对象
	ingress := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Status: nwv1.IngressStatus{},
	}
	//第一层for循环是将host组装成nwv1.IngressRule类型的对象
	//一个host对应一个ingressrule，每隔ingressrule中包含一个host和多个path
	for key, value := range data.Hosts {
		//先把host放进去
		ir := nwv1.IngressRule{
			Host: key,
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{Paths: nil},
			},
		}
		//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, httpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}
			//将每个hip对象组装成数组
			httpIngressPaths = append(httpIngressPaths, hip)
		}
		//给Paths赋值，前面置空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPaths
		//将每个ir组装成数组
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules放到ingress中
	ingress.Spec.Rules = ingressRules
	//创建ingress
	_, err = client.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("创建Ingress失败, %v\n", err))
		return errors.New(fmt.Sprintf("创建Ingress失败, %v\n", err))
	}

	return nil
}

// 删除ingress
func (i *ingress) DeleteIngress(client *kubernetes.Clientset, ingressName, namespace string) (err error) {
	err = client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Ingress失败, " + err.Error()))
		return errors.New("删除Ingress失败, " + err.Error())
	}

	return nil
}

// 更新ingress
func (i *ingress) UpdateIngress(client *kubernetes.Clientset, namespace, content string) (err error) {
	var ingress = &nwv1.Ingress{}

	err = json.Unmarshal([]byte(content), ingress)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = client.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新ingress失败, " + err.Error()))
		return errors.New("更新ingress失败, " + err.Error())
	}
	return nil
}

// workflow名字转换成ingress名字，添加-ing后缀
func getIngressName(workflowName string) (ingressName string) {
	return workflowName + "-ing"
}

func (i *ingress) toCells(std []nwv1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = ingressCell(std[i])
	}
	return cells
}

func (i *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingresss := make([]nwv1.Ingress, len(cells))
	for i := range cells {
		ingresss[i] = nwv1.Ingress(cells[i].(ingressCell))
	}

	return ingresss
}
