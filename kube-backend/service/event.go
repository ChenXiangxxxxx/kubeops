package service

import (
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"kube-backend/dao"
	"kube-backend/model"
	"time"
)

var Event event

type event struct {
}

// 获取列表
func (*event) GetList(name, cluster string, page, limit int) (*dao.Events, error) {
	data, err := dao.Event.GetList(name, cluster, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// informer监听
func (*event) WatchEventTask(cluster string) {
	//实例化
	informerFactory := informers.NewSharedInformerFactory(K8s.ClientMap[cluster], time.Minute)
	//监听资源
	informer := informerFactory.Core().V1().Events()
	//添加事件的handler
	informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				onAdd(obj, cluster)
			},
		})
	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, informer.Informer().HasSynced) {
		logger.Error("同步cache超时")
		return
	}
	<-stopCh
}

func onAdd(obj interface{}, cluster string) {
	//断言
	event := obj.(*corev1.Event)
	//判断是否重复
	_, has, err := dao.Event.HasEvent(event.InvolvedObject.Name,
		event.InvolvedObject.Kind,
		event.InvolvedObject.Namespace,
		event.Reason,
		event.CreationTimestamp.Time,
		cluster)
	if err != nil {
		return
	}
	if has {
		logger.Info(fmt.Sprintf("Event数据已存在，%s %s %s %s %v %s\n",
			event.InvolvedObject.Name,
			event.InvolvedObject.Kind,
			event.InvolvedObject.Namespace,
			event.Reason,
			event.CreationTimestamp.Time,
			cluster),
		)
	}
	data := &model.Event{
		Name:      event.InvolvedObject.Name,
		Kind:      event.InvolvedObject.Kind,
		Namespace: event.InvolvedObject.Namespace,
		Rtype:     event.Type,
		Reason:    event.Reason,
		Message:   event.Message,
		EventTime: &event.CreationTimestamp.Time,
		Cluster:   cluster,
	}
	if err := dao.Event.Add(data); err != nil {
		return
	}
}
