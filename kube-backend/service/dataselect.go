package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
)

// 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDataList   []DataCell
	dataSelectorQuery *DataSelectorQuery
}

// Datacell接口，用于各种资源list的类型转换，转换后可以使用dataSelector的方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// 定义过滤和分页的属性。过滤用name，分页用limit和page
type DataSelectorQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}
type PaginateQuery struct {
	Limit int
	Page  int
}

// 排序,实现自定义结构的排序，需要重写Len、Swap和Less的方法
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}
func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// 过滤Filter方法，比较元素的name属性。包含则返回。dataSelector从前端传参过来
func (d *dataSelector) Filter() *dataSelector {
	//Name传参为空返回所有
	if d.dataSelectorQuery.FilterQuery.Name == "" {
		return d
	}
	filteredList := []DataCell{}
	for _, value := range d.GenericDataList {
		matched := true
		objName := value.GetName()
		if !strings.Contains(objName, d.dataSelectorQuery.FilterQuery.Name) { //如果查出来的完整objName名字中不包含传入的参数的name
			matched = false
			continue
		}
		if matched {
			filteredList = append(filteredList, value)
		}
	}
	d.GenericDataList = filteredList
	return d
}

// 分页
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectorQuery.PaginateQuery.Limit
	page := d.dataSelectorQuery.PaginateQuery.Page
	//验证参数
	if limit <= 0 || page <= 0 {
		return d
	}
	//定义offset
	startIndex := limit * (page - 1)
	endIndex := limit * page
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// 定义podCell类型，实现两个方法，可进行类型转换
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func (p podCell) GetName() string {
	return p.Name
}

type deploymentCell appsv1.Deployment

func (p deploymentCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func (p deploymentCell) GetName() string {
	return p.Name
}

type serviceCell corev1.Service

func (s serviceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s serviceCell) GetName() string {
	return s.Name
}
