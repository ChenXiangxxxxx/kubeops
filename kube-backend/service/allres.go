package service

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var AllRes allRes

type allRes struct{}

// 封装map，多线程操作。
var mt sync.Mutex

func (a *allRes) GetAllNum(client *kubernetes.Clientset) (map[string]int, []error) {
	//wg是等待所有协程执行完毕。
	var wg sync.WaitGroup
	wg.Add(12)
	errs := make([]error, 0)
	//定义map[资源名]数量
	data := make(map[string]int, 0)
	go func() {
		list, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		//为什么要使用这个方法？因为有多个协程队这个map进行操作，map是线程非安全，有并发报错。因此要给map加锁。
		addMap(data, "Nodes", len(list.Items))
		wg.Done()
	}()

	go func() {
		list, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Namespaces", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Ingresses", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "PVs", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "DaemonSets", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "StatefulSets", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.BatchV1().Jobs("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Jobs", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Services", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Deployments", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Pods", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Secrets", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "ConfigMaps", len(list.Items))
		wg.Done()
	}()
	//当wg里面的计数器为0时，就取消阻塞，继续执行，非0时，则一直阻塞
	wg.Wait()
	fmt.Println("这是map:", data)
	return data, nil
}

func addMap(mp map[string]int, resource string, num int) {
	mt.Lock()
	defer mt.Unlock()
	mp[resource] = num
}
