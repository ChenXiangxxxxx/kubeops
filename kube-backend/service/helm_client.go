package service

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
)

var HelmConfig helmConfig

type helmConfig struct {
	//这种和k8s的client不一样，行不通
	//ActionConfigMap map[string]*action.Configuration
}

// 获取helm action配置
func (*helmConfig) GetAc(cluster, namespace string) (*action.Configuration, error) {
	kubeconfig, ok := K8s.KubeConfMap[cluster]
	if !ok {
		logger.Error(fmt.Sprintf("actionconfig初始化失败，cluster不存在"))
		return nil, errors.New(fmt.Sprintf("actionconfig初始化失败，cluster不存在"))
	}
	actionConfig := new(action.Configuration)
	cf := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeconfig,
		Namespace:  &namespace,
	}
	if err := actionConfig.Init(cf, namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		logger.Error(fmt.Sprintf("actionconfig初始化失败，%v\n", err))
		return nil, errors.New(fmt.Sprintf("actionconfig初始化失败，%v\n", err))
	}
	return actionConfig, nil
}
