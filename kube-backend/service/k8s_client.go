package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kube-backend/config"
)

var K8s k8s

type k8s struct {
	ClientMap   map[string]*kubernetes.Clientset
	KubeConfMap map[string]string
}

func (k *k8s) GetClient(cluster string) (*kubernetes.Clientset, error) {
	client, ok := k.ClientMap[cluster]
	if !ok {
		return nil, errors.New(fmt.Sprintf("集群%s不存在,无法获取client\n", cluster))
	}
	return client, nil
}

func (k *k8s) Init() {
	mp := make(map[string]string, 0)
	k.ClientMap = make(map[string]*kubernetes.Clientset, 0)
	if err := json.Unmarshal([]byte(config.Kubeconfigs), &mp); err != nil {
		panic(fmt.Sprintf("kubeconfigs反序列化失败", err))
	}
	k.KubeConfMap = mp

	for key, value := range mp {
		conf, err := clientcmd.BuildConfigFromFlags("", value)
		if err != nil {
			panic(fmt.Sprintf("集群%s:创建k8s配置失败 %v\n", key, err))
		}
		clientSet, err := kubernetes.NewForConfig(conf)
		if err != nil {
			panic(fmt.Sprintf("集群%s:创建k8s的client失败 %v\n", key, err))
		}
		k.ClientMap[key] = clientSet
		logger.Info(fmt.Sprintf("集群%s：创建client成功", key))
	}
}
