package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kube-backend/config"
)

var Pod pod

type pod struct{}

// 定义列表的返回类型
type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

// 获取pod列表.多集群的话选择集群会带client参数过来。
func (p *pod) GetPods(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取pod列表失败，%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取pod列表失败，%v\n", err)) //logger给自己控制看，errors.New给用户看
	}
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
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
	pods := p.fromCells(data.GenericDataList)
	return &PodsResp{
		Items: pods,
		Total: total,
	}, nil
}
func (p *pod) GetPodDetail(client *kubernetes.Clientset, podName, namespace string) (pod *corev1.Pod, err error) {
	poddetail, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取pod详情失败，%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取pod详情失败，%v\n", err))
	}
	return poddetail, nil
}

// 删除pod
func (p *pod) DeletePod(client *kubernetes.Clientset, podName, namespace string) (err error) {
	err = client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("删除pod失败，%v\n", err))
		return errors.New(fmt.Sprintf("删除pod详情失败，%v\n", err))
	}
	return nil
}

// 更新pod,需要获取content,也就是pod的整个json体
func (p *pod) UpdatePod(client *kubernetes.Clientset, namespace, content string) (err error) {
	var pods = &corev1.Pod{}
	err = json.Unmarshal([]byte(content), &pods)
	if err != nil {
		logger.Error(fmt.Sprintf("pod反序列化失败，%v\n", err))
		return errors.New(fmt.Sprintf("pod反序列化失败，%v\n", err))
	}
	_, err = client.CoreV1().Pods(namespace).Update(context.TODO(), pods, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("更新pod失败，%v\n", err))
		return errors.New(fmt.Sprintf("更新pod失败，%v\n", err))
	}
	return nil
}

// 获取container列表
func (p *pod) GetPodContainer(client *kubernetes.Clientset, podName, namespace string) (containers []string, err error) {
	podDetail, err := p.GetPodDetail(client, podName, namespace)
	if err != nil {
		return nil, err
	}
	//从pod详情中获取容器名
	for _, container := range podDetail.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

// 获取container的日志
func (p *pod) GetPodLog(client *kubernetes.Clientset, containerName, podName, namespace string) (log string, err error) {
	//设置日志的配置，容器名以及tail行数
	lineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	//获取request实例
	req := client.CoreV1().Pods(namespace).GetLogs(podName, option)
	//发起request请求，返回一个ioReadCloser类型
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(fmt.Sprintf("获取pod日志失败,%v\n", err))
		return "", errors.New(fmt.Sprintf("获取pod日志失败,%v\n", err))
	}
	defer podLogs.Close() //podLogs是一个response类型，需要close
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error(fmt.Sprintf("复制podLog失败,%v\n", err))
		return "", errors.New(fmt.Sprintf("复制podLog失败,%v\n", err))
	}
	return buf.String(), nil
}

// 定义DataCell到Pod类型转换的方法
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i]) //这里是std从podCell的结构体转位DataCell的接口类型。
	}
	return cells
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell)) //cells是接口，断言成podCell类型，再通过corev1.Pod转换。
	}
	return pods
}
